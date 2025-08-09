<script lang="ts">
  import AlgorithmSelector from './components/algorithm/AlgorithmSelector.svelte';
  import DataInput from './components/ui/DataInput.svelte';
  import VisualizationPanel from './components/visualization/VisualizationPanel.svelte';
  import BenchmarkPanel from './components/benchmark/BenchmarkPanel.svelte';
  import { selectedAlgorithm } from './stores/algorithm';
  import type { DataType, DataPattern } from './types/data';

  // 本地状态
  let inputData: any = null;
  let currentView: 'visualization' | 'benchmark' = 'visualization';

  // 事件处理
  function handleDataChange(event: CustomEvent) {
    inputData = event.detail;
    console.log('Data changed:', inputData);
  }

  function handleDataGenerate(event: CustomEvent<{ dataType: DataType; size: number; pattern: DataPattern }>) {
    const { dataType, size, pattern } = event.detail;
    console.log('Generate data:', { dataType, size, pattern });

    // 这里可以调用数据生成API
    // 暂时生成一些示例数据
    if (dataType === 'array') {
      const data = Array.from({ length: size }, (_, i) => {
        switch (pattern) {
          case 'sorted':
            return i + 1;
          case 'reversed':
            return size - i;
          case 'random':
          default:
            return Math.floor(Math.random() * 100);
        }
      });
      inputData = data;
    }
  }

  function switchView(view: 'visualization' | 'benchmark') {
    currentView = view;
  }
</script>

<main>
  <header>
    <h1>算法洞察平台</h1>
    <p>交互式算法学习和性能分析工具</p>
  </header>

  <!-- 导航标签 -->
  <div class="nav-tabs">
    <button
      class="nav-tab"
      class:active={currentView === 'visualization'}
      on:click={() => switchView('visualization')}
    >
      🎯 交互式可视化
    </button>
    <button
      class="nav-tab"
      class:active={currentView === 'benchmark'}
      on:click={() => switchView('benchmark')}
    >
      📊 性能评测对比
    </button>
  </div>

  <div class="container">
    <!-- 左侧面板 -->
    <div class="left-panel">
      <AlgorithmSelector />

      <DataInput
        bind:value={inputData}
        on:dataChange={handleDataChange}
        on:generate={handleDataGenerate}
      />

      <!-- 选中的算法信息 -->
      {#if $selectedAlgorithm}
        <div class="algorithm-info">
          <h3>选中算法: {$selectedAlgorithm.name}</h3>
          <p>{$selectedAlgorithm.description}</p>
          <div class="complexity-info">
            <span>时间复杂度: {$selectedAlgorithm.timeComplexity}</span>
            <span>空间复杂度: {$selectedAlgorithm.spaceComplexity}</span>
          </div>
        </div>
      {/if}
    </div>

    <!-- 右侧面板 -->
    <div class="right-panel">
      {#if currentView === 'visualization'}
        <VisualizationPanel {inputData} />
      {:else if currentView === 'benchmark'}
        <BenchmarkPanel />
      {/if}
    </div>
  </div>
</main>

<style>
  main {
    min-height: 100vh;
    background: var(--color-background);
    background-image: 
      radial-gradient(circle at 20% 80%, rgba(59, 130, 246, 0.03) 0%, transparent 50%),
      radial-gradient(circle at 80% 20%, rgba(139, 92, 246, 0.03) 0%, transparent 50%);
  }

  header {
    text-align: center;
    padding: var(--spacing-2xl) var(--spacing-md);
    background: var(--color-surface);
    border-bottom: 1px solid var(--color-border);
    position: relative;
    overflow: hidden;
  }

  header::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, 
      rgba(59, 130, 246, 0.05) 0%, 
      rgba(139, 92, 246, 0.05) 100%);
    pointer-events: none;
  }

  header h1 {
    margin: 0 0 var(--spacing-sm) 0;
    color: var(--color-text-primary);
    font-size: clamp(2rem, 4vw, 3rem);
    font-weight: 800;
    background: linear-gradient(135deg, var(--color-primary-600), var(--color-primary-700));
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    position: relative;
    z-index: 1;
  }

  header p {
    margin: 0;
    color: var(--color-text-secondary);
    font-size: clamp(1rem, 2vw, 1.125rem);
    font-weight: 500;
    position: relative;
    z-index: 1;
  }

  .nav-tabs {
    display: flex;
    justify-content: center;
    gap: var(--spacing-sm);
    margin: var(--spacing-xl) auto;
    padding: 0 var(--spacing-md);
    max-width: 600px;
  }

  .nav-tab {
    position: relative;
    padding: var(--spacing-md) var(--spacing-lg);
    background: var(--color-surface);
    border: 2px solid var(--color-border-light);
    border-radius: var(--radius-xl);
    font-size: 0.9375rem;
    font-weight: 600;
    cursor: pointer;
    transition: all var(--transition-normal);
    color: var(--color-text-secondary);
    backdrop-filter: blur(10px);
    flex: 1;
    max-width: 200px;
    
    /* 高级阴影效果 */
    box-shadow: 
      var(--shadow-sm),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
  }

  .nav-tab::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    padding: 2px;
    background: linear-gradient(135deg, transparent, rgba(59, 130, 246, 0.1));
    mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    mask-composite: xor;
    opacity: 0;
    transition: opacity var(--transition-normal);
  }

  .nav-tab:hover {
    border-color: var(--color-primary-300);
    color: var(--color-primary-700);
    transform: translateY(-2px);
    box-shadow: 
      var(--shadow-lg),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }

  .nav-tab:hover::before {
    opacity: 1;
  }

  .nav-tab.active {
    background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
    border-color: var(--color-primary-500);
    color: white;
    transform: translateY(-1px);
    box-shadow: 
      var(--shadow-lg),
      0 0 20px rgba(59, 130, 246, 0.3);
  }

  .nav-tab.active::before {
    opacity: 0;
  }

  .container {
    display: grid;
    grid-template-columns: 420px 1fr;
    gap: var(--spacing-xl);
    max-width: 1600px;
    margin: 0 auto;
    padding: 0 var(--spacing-md) var(--spacing-2xl);
    align-items: start;
  }

  .left-panel,
  .right-panel {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .left-panel {
    position: sticky;
    top: var(--spacing-lg);
  }

  .algorithm-info {
    padding: var(--spacing-lg);
    background: var(--color-surface);
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    transition: all var(--transition-normal);
    position: relative;
    overflow: hidden;
  }

  .algorithm-info::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: linear-gradient(90deg, var(--color-primary-500), var(--color-primary-600));
  }

  .algorithm-info:hover {
    box-shadow: var(--shadow-md);
    transform: translateY(-2px);
  }

  .algorithm-info h3 {
    margin: 0 0 var(--spacing-sm) 0;
    color: var(--color-text-primary);
    font-size: 1.125rem;
    font-weight: 700;
  }

  .algorithm-info p {
    margin: 0 0 var(--spacing-md) 0;
    color: var(--color-text-secondary);
    line-height: 1.6;
  }

  .complexity-info {
    display: flex;
    gap: var(--spacing-sm);
    font-size: 0.8125rem;
    font-weight: 500;
  }

  .complexity-info span {
    padding: var(--spacing-xs) var(--spacing-sm);
    background: var(--color-gray-100);
    color: var(--color-gray-700);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-light);
    transition: all var(--transition-fast);
  }

  .complexity-info span:hover {
    background: var(--color-primary-50);
    color: var(--color-primary-700);
    border-color: var(--color-primary-200);
  }

  /* 响应式设计 */
  @media (max-width: 1200px) {
    .container {
      grid-template-columns: 380px 1fr;
      gap: var(--spacing-lg);
    }
  }

  @media (max-width: 992px) {
    .container {
      grid-template-columns: 1fr;
      gap: var(--spacing-lg);
    }
    
    .left-panel {
      position: static;
    }
  }

  @media (max-width: 768px) {
    header {
      padding: var(--spacing-xl) var(--spacing-md);
    }

    .nav-tabs {
      flex-direction: column;
      max-width: 400px;
    }

    .nav-tab {
      max-width: none;
    }

    .container {
      padding: 0 var(--spacing-sm) var(--spacing-xl);
    }
  }

  @media (max-width: 480px) {
    .complexity-info {
      flex-direction: column;
    }
  }

  /* 平滑滚动 */
  html {
    scroll-behavior: smooth;
  }

  /* 加载动画 */
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

  .left-panel > *,
  .right-panel > * {
    animation: fadeInUp 0.6s ease-out;
  }

  .left-panel > *:nth-child(2) {
    animation-delay: 0.1s;
  }

  .left-panel > *:nth-child(3) {
    animation-delay: 0.2s;
  }
</style>
