
// Base64 to ArrayBuffer
function bufferDecode(value: string) {
  let buffer = window.atob(value);
  //let buffer = Buffer.from(value, 'base64');
  return Uint8Array.from(buffer, c => c.charCodeAt(0));
}

function bufferEncode(value: ArrayBuffer) {
  var u8Array = new Uint8Array(value) 
  let result = Array.from(u8Array);
  return window.btoa(String.fromCharCode.apply(null, result))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");;

}

export async function registerUser(username: string) {

  if (username === "") {
    alert("Please enter a username");
    return;
  }

  try {

    let beginResponse = await fetch(
      '/register/begin/' + username,
      { method: 'GET' })

    let credentialCreationOptions = await beginResponse.json()

    console.log("credentialCreationOptions", credentialCreationOptions)
    credentialCreationOptions.publicKey.challenge = bufferDecode(credentialCreationOptions.publicKey.challenge);
    credentialCreationOptions.publicKey.user.id = bufferDecode(credentialCreationOptions.publicKey.user.id);
    if (credentialCreationOptions.publicKey.excludeCredentials) {
      for (var i = 0; i < credentialCreationOptions.publicKey.excludeCredentials.length; i++) {
        credentialCreationOptions.publicKey.excludeCredentials[i].id = bufferDecode(credentialCreationOptions.publicKey.excludeCredentials[i].id);
      }
    }

    let credential = await navigator.credentials.create({
      publicKey: credentialCreationOptions.publicKey
    })
    console.log("credential", credential)
    let attestationObject = credential?.response.attestationObject;
    let clientDataJSON = credential?.response.clientDataJSON;
    let rawId = credential?.rawId;

    await fetch(
      '/register/finish/' + username,
      {
        method: 'POST',
        body: JSON.stringify({
          id: credential?.id,
          rawId: bufferEncode(rawId),
          type: credential?.type,
          response: {
            attestationObject: bufferEncode(attestationObject),
            clientDataJSON: bufferEncode(clientDataJSON),
          },
        })

      })

    console.log("successfully registered " + username + "!")
  } catch (error) {

    console.error("failed to register " + username, error)
  }
}

export async function loginUser(username: string) {

  if (username === "") {
    alert("Please enter a username");
    return;
  }

  try {

    let response = await fetch('/login/begin/' + username, { method: 'GET' });

    let credentialRequestOptions = await response.json();
    console.log("credentialRequestOptions", credentialRequestOptions)
    credentialRequestOptions.publicKey.challenge = bufferDecode(credentialRequestOptions.publicKey.challenge);
    credentialRequestOptions.publicKey.allowCredentials.forEach((listItem: any) => {
      listItem.id = bufferDecode(listItem.id)
    });

    let assertion = await navigator.credentials.get({
      publicKey: credentialRequestOptions.publicKey
    })
    console.log("assertion", assertion)
    let authData = assertion?.response.authenticatorData;
    let clientDataJSON = assertion?.response.clientDataJSON;
    let rawId = assertion?.rawId;
    let sig = assertion?.response.signature;
    let userHandle = assertion?.response.userHandle;

    let finishResponse = await fetch(
      '/login/finish/' + username,
      {
        method: 'POST', headers: {
          'Content-Type': 'application/json'
          // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: JSON.stringify({
          id: assertion?.id,
          rawId: bufferEncode(rawId),
          type: assertion?.type,
          response: {
            authenticatorData: bufferEncode(authData),
            clientDataJSON: bufferEncode(clientDataJSON),
            signature: bufferEncode(sig),
            userHandle: bufferEncode(userHandle),
          },
        })
      })

    let finishJson = await finishResponse.json();
    console.log("successfully logged in " + username + "!", finishJson)

    const urlParams = new URLSearchParams(window.location.search);
    let myParam = urlParams.get('redirect');
    console.log("redirect", window.location.host + myParam)
    if (!myParam) {
      myParam = '/login'
    }
    window.location.href = myParam
  } catch (error) {

    console.error("failed to register " + username, error)
  }
}
