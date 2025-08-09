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
              {CATEGORY_NAMES[algorithm.category as AlgorithmCategoryType]}
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
    gap: var(--spacing-lg);
    padding: var(--spacing-lg);
    background: var(--color-surface);
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    transition: all var(--transition-normal);
    position: relative;
    overflow: hidden;
  }

  .algorithm-selector::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, var(--color-primary-500), var(--color-primary-600));
  }

  .algorithm-selector:hover {
    box-shadow: var(--shadow-md);
  }

  .algorithm-selector.compact {
    padding: var(--spacing-md);
    gap: var(--spacing-md);
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--spacing-sm);
  }

  .header h3 {
    margin: 0;
    color: var(--color-text-primary);
    font-weight: 700;
    font-size: 1.125rem;
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
  }

  .header h3::before {
    content: '🧠';
    font-size: 1.25rem;
  }

  .loading {
    color: var(--color-text-muted);
    font-size: 0.8125rem;
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
  }

  .loading::before {
    content: '';
    width: 12px;
    height: 12px;
    border: 2px solid var(--color-border);
    border-top-color: var(--color-primary-500);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .search-section {
    position: relative;
  }

  .search-input {
    position: relative;
    group: search;
  }

  .search-input input {
    width: 100%;
    padding: var(--spacing-sm) 2.5rem var(--spacing-sm) var(--spacing-md);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    font-size: 0.8125rem;
    background: var(--color-surface);
    color: var(--color-text-primary);
    transition: all var(--transition-fast);
    position: relative;
  }

  .search-input input::placeholder {
    color: var(--color-text-muted);
  }

  .search-input input:focus {
    outline: none;
    border-color: var(--color-primary-500);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    background: var(--color-surface-elevated);
  }

  .search-input::before {
    content: '🔍';
    position: absolute;
    left: var(--spacing-sm);
    top: 50%;
    transform: translateY(-50%);
    color: var(--color-text-muted);
    font-size: 0.875rem;
    pointer-events: none;
    z-index: 1;
  }

  .clear-btn {
    position: absolute;
    right: var(--spacing-sm);
    top: 50%;
    transform: translateY(-50%);
    background: var(--color-gray-100);
    border: none;
    color: var(--color-text-muted);
    cursor: pointer;
    padding: var(--spacing-xs);
    border-radius: 50%;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    transition: all var(--transition-fast);
  }

  .clear-btn:hover {
    background: var(--color-gray-200);
    color: var(--color-text-secondary);
    transform: translateY(-50%) scale(1.1);
  }

  .category-section {
    border-bottom: 1px solid var(--color-border-light);
    padding-bottom: var(--spacing-md);
    margin-bottom: var(--spacing-sm);
  }

  .category-tabs {
    display: flex;
    gap: var(--spacing-xs);
    overflow-x: auto;
    padding-bottom: var(--spacing-xs);
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .category-tabs::-webkit-scrollbar {
    display: none;
  }

  .category-tab {
    padding: var(--spacing-xs) var(--spacing-md);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-xl);
    font-size: 0.8125rem;
    font-weight: 500;
    cursor: pointer;
    white-space: nowrap;
    transition: all var(--transition-fast);
    color: var(--color-text-secondary);
    position: relative;
    flex-shrink: 0;
  }

  .category-tab::before {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
    border-radius: inherit;
    opacity: 0;
    transition: opacity var(--transition-fast);
  }

  .category-tab:hover {
    background: var(--color-gray-50);
    border-color: var(--color-primary-300);
    color: var(--color-primary-700);
    transform: translateY(-1px);
  }

  .category-tab.active {
    background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
    color: white;
    border-color: var(--color-primary-500);
    box-shadow: var(--shadow-sm);
  }

  .error {
    padding: var(--spacing-md);
    background: var(--color-error-50);
    border: 1px solid rgba(239, 68, 68, 0.2);
    border-radius: var(--radius-md);
    color: var(--color-error-500);
    text-align: center;
  }

  .error button {
    margin-top: var(--spacing-sm);
    padding: var(--spacing-xs) var(--spacing-md);
    background: var(--color-error-500);
    color: white;
    border: none;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 0.8125rem;
    font-weight: 500;
    transition: all var(--transition-fast);
  }

  .error button:hover {
    background: #dc2626;
    transform: translateY(-1px);
  }

  .algorithms-list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
    max-height: 450px;
    overflow-y: auto;
    padding-right: var(--spacing-xs);
  }

  .empty {
    text-align: center;
    color: var(--color-text-muted);
    padding: var(--spacing-2xl);
    background: var(--color-gray-50);
    border-radius: var(--radius-md);
    border: 2px dashed var(--color-border);
  }

  .empty::before {
    content: '🔍';
    display: block;
    font-size: 2rem;
    margin-bottom: var(--spacing-sm);
    opacity: 0.5;
  }

  .algorithm-item {
    padding: var(--spacing-md);
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all var(--transition-normal);
    background: var(--color-surface);
    position: relative;
    overflow: hidden;
  }

  .algorithm-item::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background: var(--color-primary-500);
    transform: scaleY(0);
    transition: transform var(--transition-normal);
  }

  .algorithm-item:hover {
    border-color: var(--color-primary-300);
    box-shadow: var(--shadow-md);
    transform: translateY(-2px);
  }

  .algorithm-item:hover::before {
    transform: scaleY(1);
  }

  .algorithm-item.selected {
    border-color: var(--color-primary-500);
    background: linear-gradient(135deg, var(--color-primary-50), rgba(59, 130, 246, 0.02));
    box-shadow: 
      var(--shadow-md),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
  }

  .algorithm-item.selected::before {
    transform: scaleY(1);
    background: linear-gradient(180deg, var(--color-primary-500), var(--color-primary-600));
  }

  .algorithm-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: var(--spacing-sm);
    gap: var(--spacing-sm);
  }

  .algorithm-name {
    margin: 0;
    color: var(--color-text-primary);
    font-size: 0.9375rem;
    font-weight: 600;
    line-height: 1.3;
  }

  .algorithm-category {
    font-size: 0.6875rem;
    color: var(--color-primary-700);
    background: var(--color-primary-100);
    padding: var(--spacing-xs) var(--spacing-sm);
    border-radius: var(--radius-xl);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.025em;
    flex-shrink: 0;
    border: 1px solid var(--color-primary-200);
  }

  .algorithm-description {
    margin: 0 0 var(--spacing-md) 0;
    color: var(--color-text-secondary);
    font-size: 0.8125rem;
    line-height: 1.5;
  }

  .algorithm-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.75rem;
    gap: var(--spacing-sm);
  }

  .complexity {
    color: var(--color-text-muted);
    font-family: var(--font-family-mono);
    font-weight: 500;
  }

  .tags {
    display: flex;
    gap: var(--spacing-xs);
    flex-wrap: wrap;
  }

  .tag {
    background: linear-gradient(135deg, var(--color-success-50), rgba(16, 185, 129, 0.05));
    color: var(--color-success-500);
    padding: 2px var(--spacing-xs);
    border-radius: var(--radius-sm);
    font-size: 0.625rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.025em;
    border: 1px solid rgba(16, 185, 129, 0.2);
    white-space: nowrap;
  }

  /* 响应式设计 */
  @media (max-width: 768px) {
    .algorithm-selector {
      padding: var(--spacing-md);
    }
    
    .category-tabs {
      flex-wrap: wrap;
    }
    
    .algorithm-header {
      flex-direction: column;
      align-items: flex-start;
    }
    
    .algorithm-info {
      flex-direction: column;
      align-items: flex-start;
      gap: var(--spacing-xs);
    }
  }

  /* 加载动画 */
  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .algorithm-item {
    animation: fadeIn 0.3s ease-out;
  }

  .algorithm-item:nth-child(even) {
    animation-delay: 0.05s;
  }

  .algorithm-item:nth-child(odd) {
    animation-delay: 0.1s;
  }
</style>
