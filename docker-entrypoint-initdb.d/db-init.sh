#!/bin/bash
set -e 
DATABASE=${MONGO_INITDB_DATABASE}
USERNAME=${MONGO_INITDB_ROOT_USERNAME}
PASSWORD=${MONGO_INITDB_ROOT_PASSWORD}

echo "Adding admin user"

mongo <<EOF
use $DATABASE
db.createUser({ user: '$USERNAME',pwd:'$PASSWORD',roles:[{ role:'dbOwner', db:'admin' },{ role:'dbOwner', db:'$DATABASE' }]});
exit
EOF
 
echo "Complete"