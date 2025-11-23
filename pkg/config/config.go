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
			log.Fatalf("could not change working directory to [%s]: %s", os.Getenv("PLATO_WORKING_DIR"), color.Red("%v", err))
		}
	}

	// find git repo root and chdir into it, or else fail!
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not read current working directory: %s", color.Red("%v", err))
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
				log.Fatalf("could not get current OS user: %s", color.Red("%v", err))
			}
			pwd = path.Join(usr.HomeDir, ".plato")
			if err := os.Chdir(pwd); err != nil {
				log.Fatalf("could not change working directory to [%s]: %s", pwd, color.Red("%v", err))
			}
			break
		}

		pwd = path.Dir(pwd)
		if err := os.Chdir(pwd); err != nil {
			log.Fatalf("could not change working directory to [%s]: %s", pwd, color.Red("%v", err))
		}
	}
	log.Infof("current working directory: [%s]", color.Magenta(pwd))

	// configure Viper
	viper.SetConfigType("yaml")
	viper.SetConfigName("plato")
	viper.AddConfigPath(".")

	// automatic environment variable handling
	viper.SetEnvPrefix("PLATO")
	viper.AutomaticEnv()

	// if a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("using config file [%s]", color.Magenta(viper.ConfigFileUsed()))
	} else { // fail if no plato.yaml was found, plato insists on it!
		log.Fatalf("could not read configuration file: %s", color.Red("%v", err))
	}

	// decrypt and/or load secrets
	LoadSecrets()

	// properly re-initialize logger again, we now have the correct intended configuration values available
	log.Initialize()

	log.Infof("plato configuration loaded and ready")
}

// load secrets into config
func LoadSecrets() {
	requireSecretsYAML := true

	// check first if plato.yaml itself actually is SOPS-encrypted
	if viper.IsSet("sops.version") && viper.IsSet("sops.mac") && viper.IsSet("sops.age") {
		requireSecretsYAML = false // we don't require an additional secrets.yaml in this case
		loadSecrets(viper.ConfigFileUsed())
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not read current working directory: %s", color.Red("%v", err))
	}
	secretsFile := filepath.Join(pwd, "secrets.yaml")
	if !file.Exists(secretsFile) {
		if requireSecretsYAML {
			log.Errorf("[%s] does not exist, cannot load any additional secrets!", color.Magenta(secretsFile))
		}
		return
	}
	loadSecrets(secretsFile)
}

func loadSecrets(inputFile string) {
	// decrypt file
	decryptedSecrets, err := command.ExecOutput([]string{"sops", "-d", inputFile})
	if err != nil {
		log.Errorf("could not decrypt [%s] with SOPS: %s", color.Magenta(inputFile), color.Red("%v", err))
		return
	}
	// read in decrypted secrets
	if err := viper.MergeConfig(strings.NewReader(decryptedSecrets)); err == nil {
		log.Infof("loaded secrets from [%s]", color.Magenta(inputFile))
	} else { // fail if no secrets.yaml was found, plato insists on it!
		log.Fatalf("could not load secrets from [%s]: %s", color.Magenta(inputFile), color.Red("%v", err))
	}
}
