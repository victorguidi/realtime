<script lang="ts">
	console.log('HEre');

	type Message = {
		message: string;
	};

	let arr = [[{ sender: [] }, { receiver: [] }]];
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
		</div>
		<div class="flex flex-col bg-blue-400 w-5/6 h-full justify-around p-4">
			<div class="flex flex-col p-4 w-full h-full bg-yellow-400">
				{#each arr as item}
					<div class="flex flex-row justify-between">
						<div class="flex flex-col justify-start bg-purple-300 w-1/2">
							{#each item[0].sender as sender}
								<div>{sender.message}</div>
							{/each}
							<br />
						</div>
						<div class="flex flex-col justify-end bg-green-300 w-1/2">
							<br />
							{#each item[1].receiver as receiver}
								<div>{receiver.message}</div>
							{/each}
						</div>
					</div>
				{/each}
			</div>
			<div class="flex w-full h-1/4 bg-cyan-500">type</div>
		</div>
	</div>
</div>
