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
    background: #f8fafc;
  }

  header {
    text-align: center;
    padding: 2rem;
    background: white;
    border-bottom: 1px solid #e5e7eb;
    margin-bottom: 2rem;
  }

  header h1 {
    margin: 0 0 0.5rem 0;
    color: #1f2937;
    font-size: 2.5rem;
    font-weight: 700;
  }

  header p {
    margin: 0;
    color: #6b7280;
    font-size: 1.125rem;
  }

  .nav-tabs {
    display: flex;
    justify-content: center;
    gap: 1rem;
    margin-bottom: 2rem;
    padding: 0 1rem;
  }

  .nav-tab {
    padding: 0.75rem 1.5rem;
    background: white;
    border: 2px solid #e5e7eb;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    color: #6b7280;
  }

  .nav-tab:hover {
    border-color: #3b82f6;
    color: #3b82f6;
  }

  .nav-tab.active {
    background: #3b82f6;
    border-color: #3b82f6;
    color: white;
  }

  .container {
    display: grid;
    grid-template-columns: 400px 1fr;
    gap: 2rem;
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 1rem;
  }

  .left-panel,
  .right-panel {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }



  .algorithm-info {
    padding: 1.5rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .algorithm-info h3 {
    margin: 0 0 0.5rem 0;
    color: #1f2937;
  }

  .algorithm-info p {
    margin: 0 0 1rem 0;
    color: #4b5563;
    line-height: 1.5;
  }

  .complexity-info {
    display: flex;
    gap: 1rem;
    font-size: 0.875rem;
    color: #6b7280;
  }

  .complexity-info span {
    padding: 0.25rem 0.5rem;
    background: #f3f4f6;
    border-radius: 4px;
  }

  @media (max-width: 768px) {
    .container {
      grid-template-columns: 1fr;
      gap: 1rem;
    }

    header h1 {
      font-size: 2rem;
    }

    header p {
      font-size: 1rem;
    }
  }
</style>
