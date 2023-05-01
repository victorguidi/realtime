<script lang="ts">
	//TODO: Fetch user info and save it to the session store
	//TODO: Check why is cors not working

	async function validateUser() {
		const form = document.getElementById('form') as HTMLFormElement;
		const username = form.username.value;
		const password = form.password.value;

		// sessionStorage.setItem('token', username);
		// window.location.href = '/';

		const payload = {
			username: username,
			password: password
		};

		await fetch('https://localhost:8080/api/login', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(payload)
		})
			.then((res) => res.json())
			.then((data) => {
				if (data.error) {
					alert(data.error);
				} else {
					sessionStorage.setItem('token', data.token);
					sessionStorage.setItem('user', JSON.stringify(data.user));
					window.location.href = '/';
				}
			})
			.catch((err) => console.log(err));
	}
</script>

<div class="flex flex-col justify-center items-center h-screen">
	<div class=" flex flex-col items-center justify-around w-1/4 h-1/4 p-4">
		<h1>Login</h1>
		<div class="flex flex-col">
			<form action="" class="flex flex-col mb-5" on:submit={validateUser} id="form">
				<input
					type="text"
					name="username"
					id="username"
					placeholder="Username"
					class="p-2 m-2 border-2 border-gray-500 rounded-md"
				/>
				<input
					type="password"
					name="password"
					id="password"
					placeholder="Password"
					class="p-2 m-2 border-2 border-gray-500 rounded-md"
				/>
				<input type="submit" value="Login" class="p-2 m-2 rounded-md bg-green-800" />
			</form>
		</div>
	</div>
</div>
