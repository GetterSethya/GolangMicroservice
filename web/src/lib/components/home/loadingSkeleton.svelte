<div class="w-full h-full p-5 flex flex-col gap-5">
    {#each Array(3) as _}
        <div class="flex flex-row items-start gap-2.5">
            <div class="w-12 h-12 bg-surface-800 rounded-full animate-pulse">
                <br />
            </div>
            <div class="flex flex-col w-[calc(100%-3rem)] gap-2.5 animate-pulse">
                <div class="w-full flex flex-row justify-between">
                    <div class="bg-surface-800 rounded-lg w-1/2">
                        <br />
                    </div>
                    <div class="bg-surface-800 rounded-lg w-[10%]">
                        <br />
                    </div>
                </div>
                <div class="bg-surface-800 rounded-lg w-full">
                    <br />
                </div>
            </div>
        </div>
        <div class="w-full bg-surface-800 h-[15vh] rounded-lg animate-pulse">
            <br />
        </div>
    {/each}
</div>
