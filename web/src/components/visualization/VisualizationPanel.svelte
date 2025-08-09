<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import VisualizationCanvas from './VisualizationCanvas.svelte';
  import VisualizationControls from './VisualizationControls.svelte';
  import { 
    visualizationResult,
    visualizationLoading,
    visualizationError,
    currentStepData,
    visualizationActions,
    visualizationUtils
  } from '../../stores/visualization';
  import { selectedAlgorithm } from '../../stores/algorithm';

  // 组件属性
  export let inputData: any = null;
  export let autoExecute = false;

  // 本地状态
  let canvasWidth = 800;
  let canvasHeight = 400;
  let panelElement: HTMLElement;

  // 响应式语句
  $: if (autoExecute && $selectedAlgorithm && inputData) {
    executeVisualization();
  }

  // 生命周期
  onMount(() => {
    updateCanvasSize();
    window.addEventListener('resize', updateCanvasSize);
  });

  onDestroy(() => {
    window.removeEventListener('resize', updateCanvasSize);
    visualizationActions.clear();
  });

  // 更新画布大小
  function updateCanvasSize() {
    if (panelElement) {
      const rect = panelElement.getBoundingClientRect();
      canvasWidth = Math.min(rect.width - 40, 800);
      canvasHeight = Math.min(rect.height * 0.6, 400);
    }
  }

  // 执行可视化
  async function executeVisualization() {
    if (!$selectedAlgorithm || !inputData) {
      return;
    }

    await visualizationActions.execute($selectedAlgorithm.id, inputData);
  }

  // 重新执行
  function handleReExecute() {
    executeVisualization();
  }

  // 清除结果
  function handleClear() {
    visualizationActions.clear();
  }

  // 导出图片
  function handleExport() {
    // 这里可以实现导出功能
    console.log('Export visualization');
  }

  // 获取执行统计信息
  function getExecutionStats() {
    if (!$visualizationResult) return null;

    const stats = $visualizationResult.statistics;
    return {
      executionTime: visualizationUtils.formatExecutionTime($visualizationResult.executionTime),
      memoryUsage: visualizationUtils.formatMemoryUsage($visualizationResult.memoryUsage),
      totalSteps: $visualizationResult.totalSteps,
      comparisons: stats.comparisons,
      swaps: stats.swaps,
      moves: stats.moves,
      accesses: stats.accesses
    };
  }
</script>

<div class="visualization-panel" bind:this={panelElement}>
  <div class="panel-header">
    <h3>算法可视化</h3>
    
    <div class="header-actions">
      {#if $visualizationResult}
        <button class="action-btn" on:click={handleReExecute}>
          重新执行
        </button>
        <button class="action-btn" on:click={handleExport}>
          导出图片
        </button>
        <button class="action-btn secondary" on:click={handleClear}>
          清除
        </button>
      {:else if $selectedAlgorithm && inputData}
        <button class="action-btn primary" on:click={executeVisualization}>
          开始可视化
        </button>
      {/if}
    </div>
  </div>

  <!-- 加载状态 -->
  {#if $visualizationLoading}
    <div class="loading-state">
      <div class="loading-spinner"></div>
      <p>正在执行算法...</p>
    </div>
  {/if}

  <!-- 错误状态 -->
  {#if $visualizationError}
    <div class="error-state">
      <div class="error-icon">⚠️</div>
      <h4>执行失败</h4>
      <p>{$visualizationError}</p>
      <button class="retry-btn" on:click={executeVisualization}>
        重试
      </button>
    </div>
  {/if}

  <!-- 空状态 -->
  {#if !$visualizationResult && !$visualizationLoading && !$visualizationError}
    <div class="empty-state">
      <div class="empty-icon">🎯</div>
      <h4>准备开始可视化</h4>
      {#if !$selectedAlgorithm}
        <p>请先选择一个算法</p>
      {:else if !inputData}
        <p>请输入或生成数据</p>
      {:else}
        <p>点击"开始可视化"按钮执行算法</p>
        <button class="start-btn" on:click={executeVisualization}>
          开始可视化
        </button>
      {/if}
    </div>
  {/if}

  <!-- 可视化内容 -->
  {#if $visualizationResult && !$visualizationLoading}
    <div class="visualization-content">
      <!-- 画布 -->
      <div class="canvas-container">
        <VisualizationCanvas 
          width={canvasWidth} 
          height={canvasHeight}
        />
      </div>

      <!-- 控制器 -->
      <div class="controls-container">
        <VisualizationControls />
      </div>

      <!-- 统计信息 -->
      {#if getExecutionStats()}
        <div class="stats-container">
          <h4>执行统计</h4>
          <div class="stats-grid">
            <div class="stat-item">
              <span class="stat-label">执行时间:</span>
              <span class="stat-value">{getExecutionStats()?.executionTime}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">内存使用:</span>
              <span class="stat-value">{getExecutionStats()?.memoryUsage}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">总步骤:</span>
              <span class="stat-value">{getExecutionStats()?.totalSteps}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">比较次数:</span>
              <span class="stat-value">{getExecutionStats()?.comparisons}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">交换次数:</span>
              <span class="stat-value">{getExecutionStats()?.swaps}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">移动次数:</span>
              <span class="stat-value">{getExecutionStats()?.moves}</span>
            </div>
          </div>
        </div>
      {/if}

      <!-- 当前步骤详情 -->
      {#if $currentStepData}
        <div class="step-details">
          <h4>步骤详情</h4>
          <div class="step-content">
            <div class="step-description">
              {$currentStepData.description}
            </div>
            
            {#if $currentStepData.metadata?.phase}
              <div class="step-phase">
                阶段: {$currentStepData.metadata.phase}
              </div>
            {/if}
            
            {#if $currentStepData.metadata?.notes && $currentStepData.metadata.notes.length > 0}
              <div class="step-notes">
                <strong>备注:</strong>
                <ul>
                  {#each $currentStepData.metadata.notes as note}
                    <li>{note}</li>
                  {/each}
                </ul>
              </div>
            {/if}
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .visualization-panel {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
    padding: var(--spacing-lg);
    background: var(--color-surface);
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    min-height: 600px;
    transition: box-shadow var(--transition-normal);
    position: relative;
    overflow: hidden;
  }

  .visualization-panel::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, var(--color-primary-500), var(--color-primary-600));
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: var(--spacing-md);
    border-bottom: 1px solid var(--color-border-light);
    margin-bottom: var(--spacing-md);
  }

  .panel-header h3 {
    margin: 0;
    color: var(--color-text-primary);
    font-weight: 700;
    font-size: 1.125rem;
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
  }
  .panel-header h3::before {
    content: '🎬';
    font-size: 1.25rem;
  }

  .header-actions {
    display: flex;
    gap: var(--spacing-sm);
  }

  .action-btn {
    padding: var(--spacing-sm) var(--spacing-md);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-surface);
    color: var(--color-text-primary);
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-fast);
    box-shadow: var(--shadow-xs);
    position: relative;
  }
  .action-btn:hover {
    background: var(--color-gray-50);
    border-color: var(--color-primary-300);
    color: var(--color-primary-700);
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }
  .action-btn.primary {
    background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
    color: white;
    border-color: var(--color-primary-500);
    box-shadow: var(--shadow-sm);
  }
  .action-btn.primary:hover {
    background: linear-gradient(135deg, var(--color-primary-600), var(--color-primary-700));
    border-color: var(--color-primary-600);
  }
  .action-btn.secondary {
    color: var(--color-error-500);
    border-color: var(--color-error-50);
    background: var(--color-error-50);
  }
  .action-btn.secondary:hover {
    background: var(--color-error-50);
    border-color: var(--color-error-500);
    color: var(--color-error-500);
  }

  .loading-state,
  .error-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--spacing-2xl);
    text-align: center;
    background: var(--color-surface);
    border-radius: var(--radius-lg);
    border: 2px dashed var(--color-border-light);
    box-shadow: var(--shadow-xs);
    margin-bottom: var(--spacing-lg);
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 4px solid var(--color-gray-100);
    border-top: 4px solid var(--color-primary-500);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: var(--spacing-md);
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .error-icon,
  .empty-icon {
    font-size: 2.5rem;
    margin-bottom: var(--spacing-md);
    opacity: 0.7;
  }

  .error-state h4,
  .empty-state h4 {
    margin: 0 0 var(--spacing-xs) 0;
    color: var(--color-text-primary);
    font-size: 1.125rem;
    font-weight: 700;
  }

  .error-state p,
  .empty-state p {
    margin: 0 0 var(--spacing-md) 0;
    color: var(--color-text-secondary);
  }

  .retry-btn,
  .start-btn {
    padding: var(--spacing-md) var(--spacing-lg);
    background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
    color: white;
    border: none;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    transition: all var(--transition-fast);
    box-shadow: var(--shadow-sm);
    margin-top: var(--spacing-sm);
  }
  .retry-btn:hover,
  .start-btn:hover {
    background: linear-gradient(135deg, var(--color-primary-600), var(--color-primary-700));
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }

  .visualization-content {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
    animation: fadeInUp 0.6s ease-out;
  }

  .canvas-container {
    display: flex;
    justify-content: center;
    background: var(--color-gray-50);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-xs);
    padding: var(--spacing-md);
  }

  .stats-container,
  .step-details {
    background: var(--color-surface);
    padding: var(--spacing-lg);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-xs);
    border: 1px solid var(--color-border-light);
    margin-bottom: var(--spacing-md);
  }

  .stats-container h4,
  .step-details h4 {
    margin: 0 0 var(--spacing-sm) 0;
    color: var(--color-text-primary);
    font-size: 1rem;
    font-weight: 700;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: var(--spacing-sm);
  }

  .stat-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--spacing-sm);
    background: var(--color-gray-50);
    border-radius: var(--radius-sm);
    font-size: 0.875rem;
    font-family: var(--font-family-mono);
  }

  .stat-label {
    color: var(--color-text-muted);
    font-size: 0.875rem;
  }
  .stat-value {
    color: var(--color-primary-700);
    font-weight: 700;
    font-size: 0.875rem;
  }

  .step-content {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-xs);
  }

  .step-description {
    color: var(--color-text-primary);
    font-weight: 500;
    font-size: 0.9375rem;
  }

  .step-phase {
    color: var(--color-text-secondary);
    font-size: 0.8125rem;
  }

  .step-notes {
    color: var(--color-text-secondary);
    font-size: 0.8125rem;
  }

  .step-notes ul {
    margin: var(--spacing-xs) 0 0 var(--spacing-md);
    padding: 0;
  }

  .step-notes li {
    margin-bottom: var(--spacing-xs);
  }

  @media (max-width: 768px) {
    .visualization-panel {
      padding: var(--spacing-md);
    }
    .panel-header {
      flex-direction: column;
      gap: var(--spacing-sm);
      align-items: stretch;
    }
    .header-actions {
      justify-content: center;
    }
    .stats-grid {
      grid-template-columns: 1fr;
    }
  }

  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
