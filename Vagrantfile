# -*- mode: ruby -*-
# vi: set ft=ruby :

require 'fileutils'

unless Vagrant.has_plugin?("vagrant-vbguest")
  puts 'Installing vagrant-vbguest Plugin...'
  system('vagrant plugin install vagrant-vbguest')
end

FileUtils.mkdir_p './analyzer'

Vagrant.configure("2") do |config|
  config.ssh.insert_key = false
  config.ssh.forward_agent = true
  config.ssh.forward_x11 = true

  config.vbguest.auto_update = false
  config.vm.box_check_update = false

  config.vm.define "analyzer-vm" do |vm|
    vm.vm.hostname = "analyzer-vm"
    vm.vm.box = "bento/ubuntu-22.04"

    vm.vm.provider "virtualbox" do |vb|
        vb.name = "analyzer-vm"
        vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
        vb.memory = "8192"
        vb.cpus = 4
    end

    vm.vm.synced_folder "ssa_analysis", "/home/vagrant/ssa_analysis",
        mount_options: ["dmode=775", "fmode=775"]

    vm.vm.synced_folder "blueprint", "/home/vagrant/blueprint",
        mount_options: ["dmode=775", "fmode=775"]

    vm.vm.provision "shell", inline: <<-SHELL
        sudo apt-get update
        sudo apt-get install -y git nano wget thrift-compiler build-essential gcc libc6-dev pkg-config

        wget https://go.dev/dl/go1.22.4.linux-arm64.tar.gz
        sudo tar -C /usr/local -xzf go1.22.4.linux-arm64.tar.gz
        sudo rm -r go1.22.4.linux-arm64.tar.gz
        echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/vagrant/.profile
        echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/vagrant/.bashrc
        cd ssa_analysis/analyzer
        go mod tidy
    SHELL
  end
end
