print('Start #################################################################');

db = db.getSiblingDB('fli');
db.createUser(
    {
        user: 'root',
        pwd: 'root1234',
        roles: [{ role: 'readWrite', db: 'fli' }],
    },
);

print('END #################################################################');