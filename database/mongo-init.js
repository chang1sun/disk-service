db.createUser(
        {
            user: "chang",
            pwd: "123456",
            roles: [
                {
                    role: "readWrite",
                    db: "netdisk"
                }
            ]
        }
);
