<script lang="ts">
	//TODO: Add a router page to login
	import { writable } from 'svelte/store';

	type Message = {
		where?: string;
		from: number;
		message: string;
	};
	let message: Message;
	let id = writable<number>(0);
	let sessionId: number;
	let messages = writable<Message[]>([]);
	let socket: WebSocket;

	$: id.subscribe((val) => {
		if (val !== 0) {
			sessionId = Number(val);
			openConnection(Number(val));
		}
	});

	function openConnection(id: number): void {
		socket = new WebSocket('wss://localhost:3000/wss/login', ['1']);
		socket.addEventListener('open', function (event) {
			socket.send(
				JSON.stringify({
					sessionToken: 'tes',
					id: id
				})
			);
		});
		socket.addEventListener('message', function (event) {
			if (!event.data.includes('User authenticated')) {
				console.log(event.data);
				const msg = JSON.parse(event.data);
				if (msg.message === 'close') {
					socket.close();
					return;
				}
				message = { from: msg.from, message: msg.message } as Message;
				messages.update((msgs) => [...msgs, message]);
			}
		});
	}

	function sendMessage(): void {
		socket.send(
			JSON.stringify({
				from: sessionId,
				message: document.getElementById('sendbox')?.value as string
			})
		);
		document.getElementById('sendbox')!.value = '';
	}
</script>

<div class="flex flex-col w-screen h-screen bg-zinc-950 p-1 text-white">
	<nav class="w-full h-12">Header</nav>
	<div class="flex w-full items-center justify-around h-full">
		<div class="flex flex-col w-1/6 h-full p-4">
			<h1>Sessions</h1>
			<div class="flex p-1">
				<div>Image</div>
				<div>Session Name</div>
				<div>Users Inside</div>
			</div>
			<div class="flex">
				<h1>Add</h1>
				<h1>Create</h1>
			</div>
			<input type="number" name="id" id="id" class="text-black" />
			<button on:click={() => id.set(Number(document.getElementById('id')?.value))}>Open</button>
		</div>
		<div class="flex flex-col bg-zinc-900 w-5/6 h-full justify-around p-4 relative">
			<div class="flex flex-row p-4 w-full h-full pr-6">
				<div class="flex flex-col w-1/2 h-full items-start">
					{#each $messages as msg}
						{#if msg.from == sessionId}
							<div class="flex justify-end rounded-md p-2 h-24" />
							<!-- <br /> -->
						{/if}
						{#if msg.from != sessionId}
							<div
								class="flex items-center bg-indigo-800 min-w-[20%] max-w-prose rounded-md p-2 h-24 max-h-full"
							>
								<p class="text-white break-words">{msg.message}</p>
							</div>
							<br />
						{/if}
						{#if msg.from == sessionId}
							<div class="flex justify-end rounded-md p-2 h-24" />
							<!-- <br /> -->
						{/if}
						<!-- <div class="flex justify-end rounded-md p-2 h-24" /> -->
					{/each}
				</div>
				<div />
				<div class="flex flex-col w-1/2 h-full items-end min-w-1/4">
					{#each $messages as msg}
						{#if msg.from != sessionId}
							<div class="flex justify-end rounded-md p-2 h-24" />
							<!-- <br /> -->
						{/if}
						{#if msg.from == sessionId}
							<div
								class="flex items-center bg-indigo-800 min-w-[20%] max-w-prose rounded-md p-2 h-24 max-h-full"
							>
								<p class="text-white break-words">{msg.message}</p>
							</div>
							<br />
						{/if}
						{#if msg.from != sessionId}
							<div class="flex justify-end rounded-md p-2 h-24" />
							<!-- <br /> -->
						{/if}
					{/each}
				</div>
			</div>
			<div class="flex w-full p-2 bg-cyan-900">
				<textarea
					type="text"
					name="sendbox"
					id="sendbox"
					value=""
					class="flex w-full h-14 rounded-md text-black"
				/>
				<button on:click={sendMessage}>Send</button>
			</div>
		</div>
	</div>
</div>
