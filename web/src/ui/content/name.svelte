<script lang="ts">
    export let name: string
</script>

<span class="font-bold">{name.length > 25 ? name.substring(0, 22) + "..." : name}</span>
