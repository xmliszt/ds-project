export function parseJwt(token) {
  var base64Url = token.split(".")[1];
  var base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  var jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map(function (c) {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join("")
  );

  return JSON.parse(jsonPayload);
}

export function getRandomColor() {
  const red = Math.floor(((1 + Math.random()) * 256) / 2);
  const green = Math.floor(((1 + Math.random()) * 256) / 2);
  const blue = Math.floor(((1 + Math.random()) * 256) / 2);
  return "rgb(" + red + ", " + green + ", " + blue + ")";
}

export function deepCompareObject(objA, objB) {
  for (var k of Object.keys(objA)) {
    if (objA[k] !== objB[k]) {
      return false;
    }
  }
  return true;
}

export function deepCompareArray(arrA, arrB) {
  return JSON.stringify(arrA) === JSON.stringify(arrB);
}

export function itemExistsInArray(arr, item, key) {
  for (var e of arr) {
    if (e[key] === item) {
      return true;
    }
  }
  return false;
}
