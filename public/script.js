console.log('sanity')

const form = {
  username: '',
  password: '',
}

const usernameInput = document.getElementById('username')
const passwordInput = document.getElementById('password')
const submitButton = document.getElementById('submitBtn')

submitButton.onclick = async ev => {
  ev.preventDefault()
  console.log('click')
  try {
    // fetch on the accountHandler
    const res = await fetch(`/api/users/${ev.target.name}/accounts/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: usernameInput.value,
        password: passwordInput.value,
      }),
    })
    const data = await res.json()
    console.log(data)
  } catch (err) {
    console.log(err)
  }
}
