//TODO: change it to actual get secrets API call
export async function getSecrets(role) {
  console.log(role);
  let secrets = [];
  for (var i = 0; i < 20; i++) {
    let chars =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    // Pick characers randomly
    let str = "";
    for (let i = 0; i < 20; i++) {
      str += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    secrets.push({
      alias: str,
      value: Math.random().toString(16).substr(2, 18),
      edit: false,
      show: false,
    });
  }
  return {
    success: true,
    data: secrets,
  };
}
