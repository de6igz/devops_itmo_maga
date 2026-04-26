output "server_id" {
  value = twc_server.lab_server.id
}

output "server_name" {
  value = twc_server.lab_server.name
}

output "floating_ip" {
  value = twc_floating_ip.lab_ip.ip
}