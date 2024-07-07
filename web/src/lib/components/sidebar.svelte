<script lang="ts">
    import HomeSidebarBtn from "@lib/components/homeSidebarBtn.svelte"
    import House from "@lib/components/svg/house.svelte"
    import Person from "@lib/components/svg/person.svelte"
    import Search from "@lib/components/svg/search.svelte"
    import Power from "@lib/components/svg/power.svelte"
    import { push } from "svelte-spa-router"

    function handleLogout() {
        console.log("click")
        localStorage.removeItem("accessToken")
        localStorage.removeItem("refreshToken")

        push("/login")
    }
</script>

<section
    class="md:h-full bg-surface-900 z-40 justify-evenly h-fit gap-2.5 sm:mt-auto md:gap-5 flex flex-row md:flex-col py-5 fill-surface-400 px-5 md:px-24 md:border-e border-t border-surface-700"
>
    <HomeSidebarBtn isActive={true} href="/#/" label="Home">
        <House />
    </HomeSidebarBtn>
    <HomeSidebarBtn href="/#/profile" label="Profile">
        <Person />
    </HomeSidebarBtn>
    <HomeSidebarBtn href="/#/search" label="Search">
        <Search />
    </HomeSidebarBtn>

    <div class="mt-auto cursor-pointer">
        <button
            class="flex
            flex-row
            gap-2.5
            hover:outline
            outline-[1px]
            outline-surface-700
            px-5
            py-2.5 hover:text-primary-500 hover:fill-primary-500 w-fit rounded-lg items-center hover:bg-surface-800"
            on:click={handleLogout}
        >
            <Power />

            <span class="font-bold text-lg hidden md:block"> Logout</span>
        </button>
    </div>
</section>
