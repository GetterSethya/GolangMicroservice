<script lang="ts">
    import { cn } from "@lib/utils"

    export let profileUrl: string
    export let width = "w-[50px]"
    export let height = "h-[50px]"
</script>

<img
    src={"http://localhost/v1/image/thumbnail/" + profileUrl}
    alt=""
    class={cn(["rounded-full object-cover", width, height])}
/>
