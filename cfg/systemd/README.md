This directory contains a goplayspace.service descriptor file that can be used to install and run Go Play Space server as a service under any Linux that has systemd installed.

The descriptor file assumes that the project root directory is /var/www/goplay.space.

### Installing the service

1. Run `useradd goplayspace` to create a dedicated user to run on behalf of.
2. Copy `goplayspace.service` file to `/usr/lib/systemd/system` directory

### Starting the service

    service goplayspace start

### Stoping the service

    service goplayspace stop
