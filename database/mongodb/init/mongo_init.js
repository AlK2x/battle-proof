const DB_NAME     = process.env.MONGO_INITDB_DATABASE;
const DB_USER     = process.env.MONGODB_USER;
const DB_PASS     = process.env.MONGODB_PASS;

db = db.getSiblingDB(DB_NAME);

db.createUser({
  user: DB_USER,
  pwd: DB_PASS,
  roles: [
    { role: "readWrite", db: DB_NAME }
  ]
});