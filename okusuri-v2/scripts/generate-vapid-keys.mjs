import webpush from 'web-push';

// VAPIDキーを生成
const vapidKeys = webpush.generateVAPIDKeys();

console.log('=== VAPID Keys Generated ===');
console.log('Public Key:', vapidKeys.publicKey);
console.log('Private Key:', vapidKeys.privateKey);
console.log('\n=== Environment Variables ===');
console.log('VITE_VAPID_PUBLIC_KEY=' + vapidKeys.publicKey);
console.log('VAPID_PRIVATE_KEY=' + vapidKeys.privateKey);
console.log('\n=== Backend Environment Variables ===');
console.log('VAPID_PUBLIC_KEY=' + vapidKeys.publicKey);
console.log('VAPID_PRIVATE_KEY=' + vapidKeys.privateKey);
