package envsync

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// EnvSyncer describes some contracts to synchronize env.
type EnvSyncer interface {
	// Sync synchronizes source and target.
	// Source is the default env or the sample env.
	// Target is the actual env.
	// Both source and target are string and indicate the location of the files.
	//
	// Any values in source that aren't in target will be written to target.
	// Any values in source that are in target won't be written to target.
	Sync(source, target string) (map[string]string, error)
}

// Syncer implements EnvSyncer.
type Syncer struct {
}

// Sync implements EnvSyncer.
// Sync will read the file line by line.
// It will read the first '=' character.
// All characters prior to the first '=' character is considered as the key.
// All characters after the first '=' character until a newline character is considered as the value.
//
// e.g: FOO=bar.
// FOO is the key and bar is the value.
//
// During the synchronization process, there may be an error.
// Any key-values that have been synchronized before the error occurred is kept in target.
// Any key-values that haven't been synchronized because of an error occurred is ignored.
func (s *Syncer) Sync(source, target string) (map[string]string, error) {
	// open the source file
	sFile, err := os.Open(source)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open source file")
	}
	defer sFile.Close()

	// open the target file
	tFile, err := os.OpenFile(target, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open target file")
	}
	defer tFile.Close()

	sMap, err := godotenv.Parse(sFile)
	if err != nil {
		return nil, err
	}

	tMap, err := godotenv.Parse(tFile)
	if err != nil {
		return nil, err
	}

	addedEnv := s.additionalEnv(sMap, tMap)
	err = s.writeEnv(tFile, addedEnv)
	if err != nil {
		return nil, err
	}

	return addedEnv, nil
}

func (s *Syncer) additionalEnv(sMap, tMap map[string]string) map[string]string {
	addedEnv := make(map[string]string)
	for k, v := range sMap {
		if _, found := tMap[k]; !found {
			addedEnv[k] = v
		}
	}
	return addedEnv
}

func (s *Syncer) writeEnv(file *os.File, env map[string]string) error {
	notes := fmt.Sprintf("# Merged by envsyc at %s", time.Now())
	envContent, err := godotenv.Marshal(env)
	if err != nil {
		return errors.Wrap(err, "Error parsing data")
	}
	_, err = file.WriteString(fmt.Sprintf("%s\n%s", notes, envContent))
	if err != nil {
		return errors.Wrap(err, "Error writing to file")
	}
	return nil
}
