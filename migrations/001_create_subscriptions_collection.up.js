db.createCollection('subscriptions');
db.subscriptions.createIndex({ email: 1 }, { unique: true });
