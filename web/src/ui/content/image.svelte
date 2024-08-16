<script lang="ts">
    export let imageUrl: string
</script>

<img
    class="rounded-lg max-h-[50vh] w-full object-cover"
    src={"http://localhost/v1/image/original/" + imageUrl}
    alt=""
/>
