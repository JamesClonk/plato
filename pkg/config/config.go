package config

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/JamesClonk/plato/pkg/util/color"
	"github.com/JamesClonk/plato/pkg/util/command"
	"github.com/JamesClonk/plato/pkg/util/file"
	"github.com/JamesClonk/plato/pkg/util/log"
	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set
func InitConfig() {
	log.Initialize() // init logger already so we can use it for the code below

	// immediately chdir if PLATO_WORKING_DIR is set. Used for example to jump to "_fixtures/" for testing.
	if len(os.Getenv("PLATO_WORKING_DIR")) > 0 {
		if err := os.Chdir(os.Getenv("PLATO_WORKING_DIR")); err != nil {
			log.Fatalf("Could not change working directory to [%s]: %s", os.Getenv("PLATO_WORKING_DIR"), color.Red("%v", err))
		}
	}

	// find git repo root and chdir into it, or else fail!
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not read current working directory: %s", color.Red("%v", err))
	}
	// traverse directory path upwards, looking for the configFile in each folder
	for i := 0; i < 32; i++ { // let's go up at max 32 folders, that should be plenty :P
		if file.Exists(filepath.Join(pwd, "plato.yaml")) {
			break
		}
		if pwd == "/" || pwd == "." || pwd == "" {
			// change dir to ~/.plato as a last resort
			usr, err := user.Current()
			if err != nil {
				log.Fatalf("Could not get current OS user: %s", color.Red("%v", err))
			}
			pwd = path.Join(usr.HomeDir, ".plato")
			if err := os.Chdir(pwd); err != nil {
				log.Fatalf("Could not change working directory to [%s]: %s", pwd, color.Red("%v", err))
			}
			break
		}

		pwd = path.Dir(pwd)
		if err := os.Chdir(pwd); err != nil {
			log.Fatalf("Could not change working directory to [%s]: %s", pwd, color.Red("%v", err))
		}
	}
	log.Debugf("Current working directory: [%s]", color.Magenta(pwd))

	// configure Viper
	viper.SetConfigType("yaml")
	viper.SetConfigName("plato")
	viper.AddConfigPath(".")

	// automatic environment variable handling
	viper.SetEnvPrefix("PLATO")
	viper.AutomaticEnv()

	// if a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file [%s]", color.Magenta(viper.ConfigFileUsed()))
	} else { // fail if no plato.yaml was found, plato insists on it!
		log.Fatalf("Could not read configuration file: %s", color.Red("%v", err))
	}

	// properly re-initialize logger again, we now have the correct intended configuration values available
	log.Initialize()

	// now load secrets
	InitSecrets()

	log.Infof("plato configuration loaded")
}

// load secrets into config
func InitSecrets() {
	if !file.Exists("secrets.yaml") {
		log.Debugf("Could not detect location of [%s], cannot load any secrets!", color.Magenta("secrets.yaml"))
		return
	}

	// we should be already in the same dir as plato.yaml, secrets.yaml MUST be located here too!
	// decrypt secrets on the fly
	decryptedSecrets, err := command.ExecOutput([]string{"sops", "-d", "secrets.yaml"})
	if err != nil {
		log.Errorf("Could not decrypt [%s] via SOPS: %s", color.Magenta("secrets.yaml"), color.Red("%v", err))
		return
	}

	// read in decrypted secrets
	if err := viper.MergeConfig(strings.NewReader(decryptedSecrets)); err == nil {
		log.Debugf("Loaded secrets from [%s]", color.Magenta("secrets.yaml"))
	} else { // fail if no secrets.yaml was found, plato insists on it!
		log.Fatalf("Could not load secrets: %s", color.Red("%v", err))
	}
}
