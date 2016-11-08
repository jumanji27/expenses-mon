export default class Offline {
  constructor() {
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker
        .register('/service-workers.js')
        .then(() => {
          console.log('[SW] Registered');
        });

      navigator.serviceWorker.ready
        .then(() => {
          console.log('[SW] Ready');
        }
      );
    }
  }
}