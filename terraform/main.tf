terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0",
    }
  }
}

# Variables for configuration
variable "aws_regions" {
  description = "List of AWS regions to deploy to"
  type        = list(string)
  default     = ["us-east-1", "eu-west-1", "ap-southeast-1"]
}

variable "function_name" {
  description = "Name of the Lambda function"
  type        = string
  default     = "api-lens"
}

variable "lambda_timeout" {
  description = "Lambda function timeout in seconds"
  type        = number
  default     = 30
}

variable "lambda_memory" {
  description = "Lambda function memory in MB"
  type        = number
  default     = 512
}

# Provider configuration for each region
provider "aws" {
  region  = "us-east-1"
  alias   = "us_east_1"
  profile = "personal"
}

provider "aws" {
  region  = "eu-west-1"
  alias   = "eu_west_1"
  profile = "personal"
}

provider "aws" {
  region  = "ap-southeast-1"
  alias   = "ap_southeast_1"
  profile = "personal"
}

# Module for each region
module "lambda_us_east_1" {
  source        = "./modules/lambda"
  providers     = {
    aws = aws.us_east_1
  }
  function_name = "${var.function_name}-us-east-1"
  lambda_timeout = var.lambda_timeout
  lambda_memory  = var.lambda_memory
}

module "lambda_eu_west_1" {
  source        = "./modules/lambda"
  providers     = {
    aws = aws.eu_west_1
  }
  function_name = "${var.function_name}-eu-west-1"
  lambda_timeout = var.lambda_timeout
  lambda_memory  = var.lambda_memory
}

module "lambda_ap_southeast_1" {
  source        = "./modules/lambda"
  providers     = {
    aws = aws.ap_southeast_1
  }
  function_name = "${var.function_name}-ap-southeast-1"
  lambda_timeout = var.lambda_timeout
  lambda_memory  = var.lambda_memory
} 