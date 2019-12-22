variable "foo" {
  default = "t1.2xlarge"
}

resource "aws_instance" "foo" {
  ami           = "ami-0ff8a91507f77f867"
  instance_type = var.foo
}
