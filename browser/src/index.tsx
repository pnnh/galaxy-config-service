import ReactDOM from 'react-dom/client'
import React, { useEffect, useState } from 'react'
import {registerUser, loginUser} from './webauthn'
import './index.scss' 
  
function LoginPage () {
  useEffect(() => { 
  if (!window.PublicKeyCredential) {
    alert("Error: this browser does not support WebAuthn");
    return;
  }
  }, [])
  const [username, setUsername] = useState('larry');
  return <div>
      Username: 
  <input type="text" name="username" id="email" value={username} onChange={(event) => {
    setUsername(event.target.value)
  }} placeholder="i.e. foo@bar.com"/> 
  <button onClick={() => registerUser(username).then()}>Register</button>
  <button onClick={() => loginUser(username).then() }>Login</button>
  </div>
}


  const goTopElement2 = document.getElementById("root")
  if (goTopElement2) { 
    const root = ReactDOM.createRoot(goTopElement2)
    root.render(<LoginPage/>)
  }
 
console.debug('async fun2', location.pathname)

