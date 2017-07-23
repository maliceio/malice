//-------------------------------------------------------------------
// Malice settings
//-------------------------------------------------------------------

variable "download-url" {
    default = "https://releases.malice.io/malice/0.3.11/malice_0.3.11_linux_amd64.zip"
    description = "URL to download Malice"
}

variable "config" {
    description = "Configuration (text) for Malice"
}

variable "extra-install" {
    default = ""
    description = "Extra commands to run in the install script"
}

//-------------------------------------------------------------------
// AWS settings
//-------------------------------------------------------------------

variable "ami" {
    default = "ami-xxxxx"
    description = "AMI for Malice instances"
}

variable "availability-zones" {
    default = "us-east-1a,us-east-1b"
    description = "Availability zones for launching the Malice instances"
}

variable "elb-health-check" {
    default = "HTTP:3993/v1/sys/health"
    description = "Health check for Malice servers"
}

variable "instance_type" {
    default = "m3.medium"
    description = "Instance type for Malice instances"
}

variable "key-name" {
    default = "default"
    description = "SSH key name for Malice instances"
}

variable "nodes" {
    default = "2"
    description = "number of Malice instances"
}

variable "subnets" {
    description = "list of subnets to launch Malice within"
}

variable "vpc-id" {
    description = "VPC ID"
}
