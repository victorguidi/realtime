<script lang="ts">
	type Message = {
		from: number;
		message: string;
	};

	let message: Message;
	$: messages = [];

	function openConnection(): void {
		const socket = new WebSocket('wss://localhost:3000/wss/login', ['1']);
		socket.addEventListener('open', function (event) {
			socket.send(
				JSON.stringify({
					sessionToken: 'tes',
					id: 1
				})
			);
		});
		socket.addEventListener('message', function (event) {
			if (!event.data.includes('User authenticated')) {
				const msg = JSON.parse(event.data);
				console.log(msg.from);
				messages.push({ from: msg.from, message: msg.message } as Message);
			}
		});
	}

	// let arr = [[{ sender: [] }, { receiver: [] }]];
</script>

<div class="flex flex-col w-screen h-screen">
	<nav class="w-full h-12">Header</nav>
	<div class="flex w-full items-center justify-around h-full">
		<div class="flex flex-col bg-red-300 w-1/6 h-full p-4">
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
			<button on:click={openConnection}>Open</button>
		</div>
		<div class="flex flex-col bg-blue-400 w-5/6 h-full justify-around p-4">
			<div class="flex flex-col p-4 w-full h-full bg-yellow-400">
				<button on:click={() => console.log(messages)}>print</button>
				{#each messages as msg}
					<p>{msg.message}</p>
				{/each}
				<!-- {#each arr as item} -->
				<!-- 	<div class="flex flex-row justify-between"> -->
				<!-- 		<div class="flex flex-col justify-start bg-purple-300 w-1/2"> -->
				<!-- 			{#each item[0].sender as sender} -->
				<!-- 				<div>{sender.message.message}</div> -->
				<!-- 			{/each} -->
				<!-- 			<br /> -->
				<!-- 		</div> -->
				<!-- 		<div class="flex flex-col justify-end bg-green-300 w-1/2"> -->
				<!-- 			<br /> -->
				<!-- 			{#each item[1].receiver as receiver} -->
				<!-- 				<div>{receiver.message}</div> -->
				<!-- 			{/each} -->
				<!-- 		</div> -->
				<!-- 	</div> -->
				<!-- {/each} -->
			</div>
			<div class="flex w-full h-1/4 bg-cyan-500">type</div>
		</div>
	</div>
</div>
