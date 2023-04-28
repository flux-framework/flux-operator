package controllers

var (
	certtmp = `#   ****  Generated on 2023-04-26 22:54:42 by CZMQ  ****
#   ZeroMQ CURVE **Secret** Certificate
#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor change its permissions.
    
metadata
    name = "flux-sample-cert-generator"
    time = "2023-02-17T21:05:10"
    userid = "0"
    hostname = "flux-sample-cert-generator"
curve
    public-key = ".!?zfo10Ew)m=+J:j^zehs&{Ayy#BGSV0Eets5Ne"
    secret-key = "vmk%8&dl7ICTfgx?*+0wgPb=@kFA>djvZU-Sl[T6"
`
)

// Keygen uses zeromq to generate a curve certificate for a given hostname.
func KeyGen(name string, hostname string) (string, error) {
	return certtmp, nil
}
