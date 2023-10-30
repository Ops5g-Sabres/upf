#!/bin/bash

if [[ $EUID -ne 0 ]]; then
        echo "requires root access to run ansible playbook"
        exit 1
fi

ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i .rvn/ansible-hosts -e 'ansible_python_interpreter=/usr/bin/python3' update_all_hosts.yml

rvn configure

ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i .rvn/ansible-hosts -e 'ansible_python_interpreter=/usr/bin/python3' prepare_hosts.yml
