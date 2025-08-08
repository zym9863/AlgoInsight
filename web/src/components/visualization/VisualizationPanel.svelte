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
              <span class="stat-value">{getExecutionStats().executionTime}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">内存使用:</span>
              <span class="stat-value">{getExecutionStats().memoryUsage}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">总步骤:</span>
              <span class="stat-value">{getExecutionStats().totalSteps}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">比较次数:</span>
              <span class="stat-value">{getExecutionStats().comparisons}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">交换次数:</span>
              <span class="stat-value">{getExecutionStats().swaps}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">移动次数:</span>
              <span class="stat-value">{getExecutionStats().moves}</span>
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
    gap: 1rem;
    padding: 1rem;
    background: #f8fafc;
    border-radius: 8px;
    min-height: 600px;
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 1rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .panel-header h3 {
    margin: 0;
    color: #1f2937;
  }

  .header-actions {
    display: flex;
    gap: 0.5rem;
  }

  .action-btn {
    padding: 0.5rem 1rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    background: white;
    color: #374151;
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn:hover {
    background: #f3f4f6;
    border-color: #9ca3af;
  }

  .action-btn.primary {
    background: #3b82f6;
    color: white;
    border-color: #3b82f6;
  }

  .action-btn.primary:hover {
    background: #2563eb;
    border-color: #2563eb;
  }

  .action-btn.secondary {
    color: #dc2626;
    border-color: #fecaca;
  }

  .action-btn.secondary:hover {
    background: #fef2f2;
    border-color: #fca5a5;
  }

  .loading-state,
  .error-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    text-align: center;
    background: white;
    border-radius: 8px;
    border: 2px dashed #e5e7eb;
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f4f6;
    border-top: 4px solid #3b82f6;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .error-icon,
  .empty-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
  }

  .error-state h4,
  .empty-state h4 {
    margin: 0 0 0.5rem 0;
    color: #1f2937;
  }

  .error-state p,
  .empty-state p {
    margin: 0 0 1rem 0;
    color: #6b7280;
  }

  .retry-btn,
  .start-btn {
    padding: 0.75rem 1.5rem;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.875rem;
    cursor: pointer;
    transition: background 0.2s;
  }

  .retry-btn:hover,
  .start-btn:hover {
    background: #2563eb;
  }

  .visualization-content {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .canvas-container {
    display: flex;
    justify-content: center;
  }

  .stats-container,
  .step-details {
    background: white;
    padding: 1rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .stats-container h4,
  .step-details h4 {
    margin: 0 0 1rem 0;
    color: #1f2937;
    font-size: 1rem;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 0.5rem;
  }

  .stat-item {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem;
    background: #f9fafb;
    border-radius: 4px;
  }

  .stat-label {
    color: #6b7280;
    font-size: 0.875rem;
  }

  .stat-value {
    color: #1f2937;
    font-weight: 500;
    font-size: 0.875rem;
  }

  .step-content {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .step-description {
    color: #1f2937;
    font-weight: 500;
  }

  .step-phase {
    color: #6b7280;
    font-size: 0.875rem;
  }

  .step-notes {
    color: #4b5563;
    font-size: 0.875rem;
  }

  .step-notes ul {
    margin: 0.25rem 0 0 1rem;
    padding: 0;
  }

  .step-notes li {
    margin-bottom: 0.25rem;
  }

  @media (max-width: 768px) {
    .panel-header {
      flex-direction: column;
      gap: 1rem;
      align-items: stretch;
    }

    .header-actions {
      justify-content: center;
    }

    .stats-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
