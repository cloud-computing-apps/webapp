variable "aws_zone" {
  description = "The AWS Zone"
  type        = string
  default     = "us-east-1a"
}

variable "region" {
  description = "The Cloud Platform Region"
  type        = string
  default     = "us-east-1"
}

variable "profile" {
  description = "The Cloud Platform Profile"
  type        = string
  default     = "dev"
}

variable "cloud_provider" {
  description = "The Cloud Provider type"
  type        = string
  default     = "aws"
}

variable "instance_type" {
  description = "The EC2 Instance Type"
  type        = string
  default     = "t2.micro"
}

variable "root_device_type" {
  description = "The EC2 Root Device Type"
  type        = string
  default     = "ebs"
}

variable "virt_type" {
  description = "The EC2 Virtualization Type"
  type        = string
  default     = "hvm"
}

variable "image_name" {
  description = "The Image Name"
  type        = string
  default     = "webapp-image"
}

variable "source_image" {
  description = "The Source Image"
  type        = string
  default     = "ubuntu-2404-noble-amd64-v20250214"
}

variable "aws_ami_name" {
  description = "The AWS Amazon Machine Image Name"
  type        = string
  default     = "ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-20250115"
}

variable "aws_owner" {
  description = "The AWS Owner"
  type        = string
  default     = "099720109477"
}

variable "ssh_username" {
  description = "The username for ssh"
  type        = string
  default     = "ubuntu"
}

variable "device_name" {
  description = "The device name"
  type        = string
  default     = "/dev/sda1"
}

variable "tags" {
  type = map(string)
  default = {
    Name        = "webapp-image"
    CreatedBy   = "packer"
    Environment = "private"
  }
}
