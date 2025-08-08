<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    algorithms, 
    selectedAlgorithm, 
    selectedCategory,
    filteredAlgorithms,
    searchQuery,
    algorithmsLoading,
    algorithmsError,
    algorithmActions 
  } from '../../stores/algorithm';
  import { ALGORITHM_CATEGORIES, CATEGORY_NAMES } from '../../types/algorithm';
  import type { Algorithm, AlgorithmCategoryType } from '../../types/algorithm';

  // 组件属性
  export let showSearch = true;
  export let showCategories = true;
  export let compact = false;

  // 本地状态
  let searchInput = '';

  // 响应式语句
  $: {
    if (searchInput !== $searchQuery) {
      algorithmActions.setSearchQuery(searchInput);
    }
  }

  // 生命周期
  onMount(() => {
    algorithmActions.loadAlgorithms();
  });

  // 事件处理
  function handleAlgorithmSelect(algorithm: Algorithm) {
    algorithmActions.selectAlgorithm(algorithm);
  }

  function handleCategorySelect(category: AlgorithmCategoryType | null) {
    selectedCategory.set(category);
  }

  function handleSearch() {
    algorithmActions.searchAlgorithms(searchInput);
  }

  function clearSearch() {
    searchInput = '';
    algorithmActions.clearSearch();
  }

  // 获取算法复杂度显示文本
  function getComplexityText(algorithm: Algorithm): string {
    return `时间: ${algorithm.timeComplexity}, 空间: ${algorithm.spaceComplexity}`;
  }

  // 获取算法特性标签
  function getAlgorithmTags(algorithm: Algorithm): string[] {
    const tags = [];
    if (algorithm.stable) tags.push('稳定');
    if (algorithm.inPlace) tags.push('原地');
    if (algorithm.adaptive) tags.push('自适应');
    return tags;
  }
</script>

<div class="algorithm-selector" class:compact>
  <!-- 标题 -->
  <div class="header">
    <h3>算法选择</h3>
    {#if $algorithmsLoading}
      <div class="loading">加载中...</div>
    {/if}
  </div>

  <!-- 搜索框 -->
  {#if showSearch}
    <div class="search-section">
      <div class="search-input">
        <input
          type="text"
          placeholder="搜索算法..."
          bind:value={searchInput}
          on:input={handleSearch}
        />
        {#if searchInput}
          <button class="clear-btn" on:click={clearSearch}>
            ✕
          </button>
        {/if}
      </div>
    </div>
  {/if}

  <!-- 类别过滤 -->
  {#if showCategories}
    <div class="category-section">
      <div class="category-tabs">
        <button
          class="category-tab"
          class:active={!$selectedCategory}
          on:click={() => handleCategorySelect(null)}
        >
          全部
        </button>
        {#each Object.values(ALGORITHM_CATEGORIES) as category}
          <button
            class="category-tab"
            class:active={$selectedCategory === category}
            on:click={() => handleCategorySelect(category)}
          >
            {CATEGORY_NAMES[category]}
          </button>
        {/each}
      </div>
    </div>
  {/if}

  <!-- 错误信息 -->
  {#if $algorithmsError}
    <div class="error">
      <p>加载失败: {$algorithmsError}</p>
      <button on:click={() => algorithmActions.loadAlgorithms()}>
        重试
      </button>
    </div>
  {/if}

  <!-- 算法列表 -->
  <div class="algorithms-list">
    {#if $filteredAlgorithms.length === 0 && !$algorithmsLoading}
      <div class="empty">
        {#if $searchQuery}
          没有找到匹配的算法
        {:else}
          暂无算法数据
        {/if}
      </div>
    {:else}
      {#each $filteredAlgorithms as algorithm (algorithm.id)}
        <div
          class="algorithm-item"
          class:selected={$selectedAlgorithm?.id === algorithm.id}
          on:click={() => handleAlgorithmSelect(algorithm)}
          role="button"
          tabindex="0"
          on:keydown={(e) => e.key === 'Enter' && handleAlgorithmSelect(algorithm)}
        >
          <div class="algorithm-header">
            <h4 class="algorithm-name">{algorithm.name}</h4>
            <div class="algorithm-category">
              {CATEGORY_NAMES[algorithm.category]}
            </div>
          </div>
          
          <p class="algorithm-description">
            {algorithm.description}
          </p>
          
          <div class="algorithm-info">
            <div class="complexity">
              {getComplexityText(algorithm)}
            </div>
            
            {#if getAlgorithmTags(algorithm).length > 0}
              <div class="tags">
                {#each getAlgorithmTags(algorithm) as tag}
                  <span class="tag">{tag}</span>
                {/each}
              </div>
            {/if}
          </div>
        </div>
      {/each}
    {/if}
  </div>
</div>

<style>
  .algorithm-selector {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .algorithm-selector.compact {
    padding: 0.5rem;
    gap: 0.5rem;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .header h3 {
    margin: 0;
    color: #1f2937;
  }

  .loading {
    color: #6b7280;
    font-size: 0.875rem;
  }

  .search-section {
    position: relative;
  }

  .search-input {
    position: relative;
  }

  .search-input input {
    width: 100%;
    padding: 0.5rem 2rem 0.5rem 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.875rem;
  }

  .search-input input:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .clear-btn {
    position: absolute;
    right: 0.5rem;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    color: #6b7280;
    cursor: pointer;
    padding: 0.25rem;
  }

  .clear-btn:hover {
    color: #374151;
  }

  .category-section {
    border-bottom: 1px solid #e5e7eb;
  }

  .category-tabs {
    display: flex;
    gap: 0.5rem;
    overflow-x: auto;
    padding-bottom: 0.5rem;
  }

  .category-tab {
    padding: 0.5rem 1rem;
    background: none;
    border: 1px solid #d1d5db;
    border-radius: 20px;
    font-size: 0.875rem;
    cursor: pointer;
    white-space: nowrap;
    transition: all 0.2s;
  }

  .category-tab:hover {
    background: #f3f4f6;
  }

  .category-tab.active {
    background: #3b82f6;
    color: white;
    border-color: #3b82f6;
  }

  .error {
    padding: 1rem;
    background: #fef2f2;
    border: 1px solid #fecaca;
    border-radius: 6px;
    color: #dc2626;
  }

  .error button {
    margin-top: 0.5rem;
    padding: 0.25rem 0.5rem;
    background: #dc2626;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .algorithms-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    max-height: 400px;
    overflow-y: auto;
  }

  .empty {
    text-align: center;
    color: #6b7280;
    padding: 2rem;
  }

  .algorithm-item {
    padding: 1rem;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .algorithm-item:hover {
    border-color: #3b82f6;
    box-shadow: 0 2px 4px rgba(59, 130, 246, 0.1);
  }

  .algorithm-item.selected {
    border-color: #3b82f6;
    background: #eff6ff;
  }

  .algorithm-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .algorithm-name {
    margin: 0;
    color: #1f2937;
    font-size: 1rem;
  }

  .algorithm-category {
    font-size: 0.75rem;
    color: #6b7280;
    background: #f3f4f6;
    padding: 0.25rem 0.5rem;
    border-radius: 12px;
  }

  .algorithm-description {
    margin: 0 0 0.5rem 0;
    color: #4b5563;
    font-size: 0.875rem;
    line-height: 1.4;
  }

  .algorithm-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.75rem;
  }

  .complexity {
    color: #6b7280;
  }

  .tags {
    display: flex;
    gap: 0.25rem;
  }

  .tag {
    background: #dbeafe;
    color: #1e40af;
    padding: 0.125rem 0.375rem;
    border-radius: 8px;
    font-size: 0.625rem;
  }
</style>
