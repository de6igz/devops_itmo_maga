variable "twc_token" {
  type        = string
  sensitive   = true
  description = "API token Timeweb Cloud"
}

variable "server_name" {
  type        = string
  description = "Имя виртуальной машины"
  default     = "lab2-server"
}

variable "availability_zone" {
  type        = string
  description = "Зона доступности Timeweb"
  default     = "spb-3"
}

variable "location" {
  type        = string
  description = "Локация Timeweb Cloud"
  default     = "ru-1"
}

variable "preset_type" {
  type        = string
  description = "Тип конфигуратора"
  default     = "premium"
}

variable "os_name" {
  type        = string
  description = "Название ОС"
  default     = "ubuntu"
}

variable "os_version" {
  type        = string
  description = "Версия ОС"
  default     = "20.04"
}

variable "cpu" {
  type        = number
  description = "Количество vCPU"
  default     = 1
}

variable "ram" {
  type        = number
  description = "RAM в МБ"
  default     = 1024
}

variable "disk" {
  type        = number
  description = "Диск в МБ"
  default     = 10240
}

variable "ssh_public_key_path" {
  type        = string
  description = "Путь до публичного SSH-ключа"
  default     = "~/.ssh/id_rsa.pub"
}