<script lang="ts" module>
	import AudioWaveform from 'lucide-svelte/icons/audio-waveform';
	import BookOpen from 'lucide-svelte/icons/book-open';
	import Bot from 'lucide-svelte/icons/bot';
	import ChartPie from 'lucide-svelte/icons/chart-pie';
	import Command from 'lucide-svelte/icons/command';
	import Frame from 'lucide-svelte/icons/frame';
	import GalleryVerticalEnd from 'lucide-svelte/icons/gallery-vertical-end';
	import Map from 'lucide-svelte/icons/map';
	import Settings2 from 'lucide-svelte/icons/settings-2';
	import SquareTerminal from 'lucide-svelte/icons/square-terminal';
	import HomeIcon from 'lucide-svelte/icons/home';

	const data = {
		navMain: [
			{
				title: 'Home',
				url: '#',
				icon: HomeIcon,
				isActive: true,
				items: [
					{
						title: 'For You',
						url: '#'
					},
					{
						title: 'Following',
						url: '#'
					},
				]
			},
			// {
			// 	title: 'Models',
			// 	url: '#',
			// 	icon: Bot,
			// 	items: [
			// 		{
			// 			title: 'Genesis',
			// 			url: '#'
			// 		},
			// 		{
			// 			title: 'Explorer',
			// 			url: '#'
			// 		},
			// 		{
			// 			title: 'Quantum',
			// 			url: '#'
			// 		}
			// 	]
			// },
			// {
			// 	title: 'Documentation',
			// 	url: '#',
			// 	icon: BookOpen,
			// 	items: [
			// 		{
			// 			title: 'Introduction',
			// 			url: '#'
			// 		},
			// 		{
			// 			title: 'Get Started',
			// 			url: '#'
			// 		},
			// 		{
			// 			title: 'Tutorials',
			// 			url: '#'
			// 		},
			// 		{
			// 			title: 'Changelog',
			// 			url: '#'
			// 		}
			// 	]
			// },
			// {
			// 	title: 'Settings',
			// 	url: '#',
			// 	icon: Settings2,
			// 	items: [
			// 		{
			// 			title: 'General',
			// 			url: '#'
			// 		},
			// 		{
			// 			title: 'Team',
			// 			url: '#'
			// 		},
			// 		{
			// 			title: 'Billing',
			// 			url: '#'
			// 		},
			// 		{
			// 			title: 'Limits',
			// 			url: '#'
			// 		}
			// 	]
			// }
		],
		// projects: [
		// 	{
		// 		name: 'Design Engineering',
		// 		url: '#',
		// 		icon: Frame
		// 	},
		// 	{
		// 		name: 'Sales & Marketing',
		// 		url: '#',
		// 		icon: ChartPie
		// 	},
		// 	{
		// 		name: 'Travel',
		// 		url: '#',
		// 		icon: Map
		// 	}
		// ]
	};
</script>

<script lang="ts">
	import NavMain from '$lib/components/nav-main.svelte';
	import NavProjects from '$lib/components/nav-projects.svelte';
	import NavUser from '$lib/components/nav-user.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import { getLocalDataCtx } from '../../routes/(protected)/+layout.svelte';

	let {
		ref = $bindable(null),
		collapsible = 'icon',
		...restProps
	}: ComponentProps<typeof Sidebar.Root> = $props();

	const localData = getLocalDataCtx();
</script>

<Sidebar.Root bind:ref {collapsible} {...restProps}>
	<!-- <Sidebar.Header> -->
	<!-- 	<TeamSwitcher teams={data.teams} /> -->
	<!-- </Sidebar.Header> -->
	<Sidebar.Content>
		<NavMain items={data.navMain} />
		<!-- <NavProjects projects={data.projects} /> -->
	</Sidebar.Content>
	<Sidebar.Footer>
		{#if $localData.user}
			<NavUser user={$localData.user} />
		{/if}
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
