package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/plugins"
)

func cmdEnablePlugin() {
	// ~ ❯❯❯ apm enable go-plus
	// Enabled:
	//   go-plus
	// ✓
}

func cmdDisablePlugin() {
	// ~ ❯❯❯ apm disable go-plus
	// Disabled:
	//   go-plus
	// ✓
}

func cmdGoToPluginHome() {

}

func cmdPluginInit() {
	// ~ ❯❯❯ apm init
	// You must specify either --package, --theme or --language to `apm init`
}

func cmdLogin() {
	// ~ ❯❯❯ apm login
	// Welcome to Atom!

	// Before you can publish packages, you'll need an API token.

	// Visit your account page on Atom.io https://atom.io/account,
	// copy the token and paste it below when prompted.

	// Press [Enter] to open your account page on Atom.io.
}

func cmdPublishPlugin() {
	// ~ ❯❯❯ apm publish
	// No package.json file found at /Users/blacktop/package.json
}

func cmdUnpublishPlugin() {
	// ~ ❯❯❯ apm unpublish
	// Are you sure you want to unpublish ALL VERSIONS of 'plugin'? This will remove it from the apm registry,
	// including download counts and stars, and this action is irreversible. (no)

}

func cmdListPlugins(all bool, detail bool) error {
	if all {
		plugins.ListAllPlugins(detail)
	} else {
		plugins.ListEnabledPlugins(detail)
	}

	// TODO: Add ability to list malice plugins not installed

	// docker := client.NewDockerClient()
	// err := docker.SearchImages("malice")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// enabled := plugins.GetEnabledPlugins()
	// fmt.Println(enabled)
	return nil
}

func cmdShowOutdatedPlugins() {
	// ~ ❯❯❯ apm outdated
	// Package Updates Available (1)
	// └── MagicPython 0.5.15 -> 0.5.16
}

func cmdPluginSearch(name string) {
	// ~ ❯❯❯ apm search go-plus
	// Search Results For 'go-plus' (29)
	// ├── go-plus Makes working with go in Atom awesome. (213713 downloads, 372 stars)
	// ├── browser-plus Browser Plus (70140 downloads, 144 stars)
	// <SNIP>
	// ├── symbols-plus Symbol browser and (project-wide) jump (replaces built-in symbols-view). (707 downloads, 3 stars)
	// └── notepad-plus-plus-syntax Notepad++ default theme syntax (596 downloads, 2 stars)

	// Use `apm install` to install them or visit http://atom.io/packages to read more about them.
}

func cmdShowPlugin() {
	// ~ ❯❯❯ apm show go-plus                                                                                                                                                         ⏎
	// go-plus
	// ├── 4.3.2
	// ├── https://github.com/joefitzgerald/go-plus
	// ├── Makes working with go in Atom awesome.
	// ├── 213622 downloads
	// └── 372 stars

	// Run `apm install go-plus` to install this package.
}

func cmdInstallPlugin(name string) error {
	// 	~ ❯❯❯ apm install -h

	// Usage: apm install [<package_name>...]
	//        apm install <package_name>@<package_version>
	//        apm install <git_remote>
	//        apm install <github_username>/<github_project>
	//        apm install --packages-file my-packages.txt
	//        apm i (with any of the previous argument usage)

	// Install the given Atom package to ~/.atom/packages/<package_name>.

	// If no package name is given then all the dependencies in the package.json
	// file are installed to the node_modules folder in the current working
	// directory.

	// A packages file can be specified that is a newline separated list of
	// package names to install with optional versions using the
	// `package-name@version` syntax.

	// Options:
	//   --check           Check that native build tools are installed                            [boolean]
	//   --verbose         Show verbose debug information                        [boolean] [default: false]
	//   --packages-file   A text file containing the packages to install                          [string]
	//   --production      Do not install dev dependencies                                        [boolean]
	//   -c, --compatible  Only install packages/themes compatible with this Atom version          [string]
	//   -h, --help        Print this usage message
	//   -s, --silent      Set the npm log level to silent                                        [boolean]
	//   -q, --quiet       Set the npm log level to warn                                          [boolean]

	//   Prefix an option with `no-` to set it to false such as --no-color to disable
	//   colored output.
	// ==============================================================================
	// ~ ❯❯❯ apm i                                                                                                                                                                    ⏎
	// Installing modules ✓
	testPlugin := plugins.Plugin{
		Name:        name,
		Enabled:     true,
		Category:    "test",
		Description: "This is a test plugin",
		Image:       "blacktop/test",
		Mime:        "image/png",
	}
	return plugins.InstallPlugin(&testPlugin)
}

func cmdRemovePlugin(name string) error {
	return plugins.DeletePlugin(name)
}

func cmdUpdatePlugin(name string, all bool, source bool) error {
	docker := client.NewDockerClient()
	if all {
		plugins.UpdateEnabledPlugins(docker)
	} else {
		if name == "" {
			log.Error("Please enter a valid plugin name.")
			os.Exit(1)
		}
		if source {
			plugins.GetPluginByName(name).UpdatePluginFromRepository(docker)
		} else {
			plugins.GetPluginByName(name).UpdatePlugin(docker)
		}
	}
	return nil
}

func cmdTestPlugin() {
	// 	~ ❯❯❯ apm test go-plus
	// [42497:1008/100044:WARNING:resource_bundle.cc(311)] locale_file_path.empty() for locale English
	// [warn] kq_init: detected broken kqueue; not using.: Undefined error: 0
	// [42497:1008/100044:WARNING:dns_config_service_posix.cc(146)] dns_config has unhandled options!
	// [warn] kq_init: detected broken kqueue; not using.: Undefined error: 0
	// [warn] kq_init: detected broken kqueue; not using.: Undefined error: 0
	// [warn] kq_init: detected broken kqueue; not using.: Undefined error: 0
	// [warn] kq_init: detected broken kqueue; not using.: Undefined error: 0
	// [42501:1008/100045:WARNING:resource_bundle.cc(311)] locale_file_path.empty() for locale English
	// [warn] kq_init: detected broken kqueue; not using.: Undefined error: 0
	// [42497:1008/100046:INFO:CONSOLE(84)] "Error: Cannot find module '/Users/blacktop/spec'
	//     at Module._resolveFilename (module.js:339:15)
	//     at Function.Module._resolveFilename (/Applications/Atom.app/Contents/Resources/app.asar/src/module-cache.js:383:52)
	//     at Function.Module._load (module.js:290:25)
	//     at Module.require (module.js:367:17)
	//     at require (/Applications/Atom.app/Contents/Resources/app.asar/src/native-compile-cache.js:50:27)
	//     at requireSpecs (/Applications/Atom.app/Contents/Resources/app.asar/spec/jasmine-test-runner.js:101:7)
	//     at module.exports (/Applications/Atom.app/Contents/Resources/app.asar/spec/jasmine-test-runner.js:49:7)
	//     at module.exports (/Applications/Atom.app/Contents/Resources/app.asar/src/initialize-test-window.js:68:17)
	//     at setupWindow (file:///Applications/Atom.app/Contents/Resources/app.asar/static/index.js:82:12)
	//     at window.onload (file:///Applications/Atom.app/Contents/Resources/app.asar/static/index.js:41:9)", source: /Applications/Atom.app/Contents/Resources/app.asar/src/initialize-test-window.js (84)
	// [42498:1008/100046:WARNING:channel.cc(359)] RawChannel write error
	// Tests failed
}
