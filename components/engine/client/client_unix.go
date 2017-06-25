// +build linux freebsd solaris openbsd darwin

package client

// DefaultMaliceHost defines os specific default if MALICE_HOST is unset
const DefaultMaliceHost = "unix:///var/run/malice.sock"
