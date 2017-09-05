resource "template_file" "install" {
    template = "${file("${path.module}/scripts/install.sh.tpl")}"

    vars {
        download_url  = "${var.download-url}"
        config        = "${var.config}"
        extra-install = "${var.extra-install}"
    }
}

// We launch Malice into an ASG so that it can properly bring them up for us.
resource "aws_autoscaling_group" "malice" {
    name = "malice - ${aws_launch_configuration.malice.name}"
    launch_configuration = "${aws_launch_configuration.malice.name}"
    availability_zones = ["${split(",", var.availability-zones)}"]
    min_size = "${var.nodes}"
    max_size = "${var.nodes}"
    desired_capacity = "${var.nodes}"
    health_check_grace_period = 15
    health_check_type = "EC2"
    vpc_zone_identifier = ["${split(",", var.subnets)}"]
    load_balancers = ["${aws_elb.malice.id}"]

    tag {
        key = "Name"
        value = "malice"
        propagate_at_launch = true
    }
}

resource "aws_launch_configuration" "malice" {
    image_id = "${var.ami}"
    instance_type = "${var.instance_type}"
    key_name = "${var.key-name}"
    security_groups = ["${aws_security_group.malice.id}"]
    user_data = "${template_file.install.rendered}"
}

// Security group for Malice allows SSH and HTTP access (via "tcp" in
// case TLS is used)
resource "aws_security_group" "malice" {
    name = "malice"
    description = "Malice servers"
    vpc_id = "${var.vpc-id}"
}

resource "aws_security_group_rule" "malice-ssh" {
    security_group_id = "${aws_security_group.malice.id}"
    type = "ingress"
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
}

// This rule allows Malice HTTP API access to individual nodes, since each will
// need to be addressed individually for unsealing.
resource "aws_security_group_rule" "malice-http-api" {
    security_group_id = "${aws_security_group.malice.id}"
    type = "ingress"
    from_port = 8200
    to_port = 8200
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "malice-egress" {
    security_group_id = "${aws_security_group.malice.id}"
    type = "egress"
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
}

// Launch the ELB that is serving Malice. This has proper health checks
// to only serve healthy, unsealed Vaults.
resource "aws_elb" "malice" {
    name = "malice"
    connection_draining = true
    connection_draining_timeout = 400
    internal = true
    subnets = ["${split(",", var.subnets)}"]
    security_groups = ["${aws_security_group.elb.id}"]

    listener {
        instance_port = 8200
        instance_protocol = "tcp"
        lb_port = 80
        lb_protocol = "tcp"
    }

    listener {
        instance_port = 8200
        instance_protocol = "tcp"
        lb_port = 443
        lb_protocol = "tcp"
    }

    health_check {
        healthy_threshold = 2
        unhealthy_threshold = 3
        timeout = 5
        target = "${var.elb-health-check}"
        interval = 15
    }
}

resource "aws_security_group" "elb" {
    name = "malice-elb"
    description = "Malice ELB"
    vpc_id = "${var.vpc-id}"
}

resource "aws_security_group_rule" "malice-elb-http" {
    security_group_id = "${aws_security_group.elb.id}"
    type = "ingress"
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "malice-elb-https" {
    security_group_id = "${aws_security_group.elb.id}"
    type = "ingress"
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "malice-elb-egress" {
    security_group_id = "${aws_security_group.elb.id}"
    type = "egress"
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
}
