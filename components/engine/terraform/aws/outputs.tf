output "address" {
    value = "${aws_elb.malice.dns_name}"
}

// Can be used to add additional SG rules to malice instances.
output "vault_security_group" {
    value = "${aws_security_group.malice.id}"
}

// Can be used to add additional SG rules to the malice ELB.
output "elb_security_group" {
    value = "${aws_security_group.elb.id}"
}
