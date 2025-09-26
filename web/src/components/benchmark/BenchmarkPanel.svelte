<script lang="ts">
  import { onMount } from 'svelte';
  import { selectedAlgorithm, algorithms, algorithmActions } from '../../stores/algorithm';
  import type { Algorithm } from '../../types/algorithm';
  import type { DataType, DataPattern } from '../../types/data';
  import { DATA_TYPES, DATA_PATTERNS, DATA_TYPE_NAMES, PATTERN_NAMES } from '../../types/data';

  // 组件状态
  let selectedAlgorithms: string[] = [];
  let dataSizes: number[] = [100, 500, 1000, 5000];
  let dataType: DataType = DATA_TYPES.ARRAY;
  let dataPattern: DataPattern = DATA_PATTERNS.RANDOM;
  let testCount = 3;
  let isRunning = false;
  let results: any[] = [];

  const idPrefix = `benchmark-${Math.random().toString(36).slice(2)}`;
  const inputIds = {
    dataType: `${idPrefix}-data-type`,
    dataPattern: `${idPrefix}-data-pattern`,
    testCount: `${idPrefix}-test-count`
  } as const;

  // 响应式语句
  $: availableAlgorithms = $algorithms.filter(alg => 
    alg.category === 'sorting' || alg.category === 'searching' || alg.category === 'graph'
  );

  // 生命周期
  onMount(() => {
    algorithmActions.loadAlgorithms();
  });

  // 事件处理
  function toggleAlgorithmSelection(algorithmId: string) {
    if (selectedAlgorithms.includes(algorithmId)) {
      selectedAlgorithms = selectedAlgorithms.filter(id => id !== algorithmId);
    } else {
      selectedAlgorithms = [...selectedAlgorithms, algorithmId];
    }
  }

  function addDataSize() {
    const newSize = prompt('请输入数据规模:');
    if (newSize && !isNaN(Number(newSize))) {
      dataSizes = [...dataSizes, Number(newSize)].sort((a, b) => a - b);
    }
  }

  function removeDataSize(index: number) {
    dataSizes = dataSizes.filter((_, i) => i !== index);
  }

  async function runBenchmark() {
    if (selectedAlgorithms.length === 0) {
      alert('请至少选择一个算法');
      return;
    }

    if (dataSizes.length === 0) {
      alert('请至少添加一个数据规模');
      return;
    }

    isRunning = true;
    results = [];

    try {
      // 模拟性能测试
      for (const algorithmId of selectedAlgorithms) {
        const algorithm = availableAlgorithms.find(alg => alg.id === algorithmId);
        if (!algorithm) continue;

        for (const size of dataSizes) {
          // 模拟测试结果
          const avgTime = Math.random() * size * Math.log(size) / 1000;
          const result = {
            algorithmId,
            algorithmName: algorithm.name,
            dataSize: size,
            avgTime: avgTime.toFixed(2),
            minTime: (avgTime * 0.8).toFixed(2),
            maxTime: (avgTime * 1.2).toFixed(2),
            memoryUsage: (size * 4 + Math.random() * 1000).toFixed(0),
            operations: Math.floor(size * Math.log(size) + Math.random() * size)
          };
          results = [...results, result];
          
          // 模拟延迟
          await new Promise(resolve => setTimeout(resolve, 500));
        }
      }
    } catch (error) {
      console.error('Benchmark failed:', error);
      alert('性能测试失败');
    } finally {
      isRunning = false;
    }
  }

  function clearResults() {
    results = [];
  }

  function exportResults() {
    const csv = generateCSV(results);
    downloadCSV(csv, 'benchmark-results.csv');
  }

  function generateCSV(data: any[]): string {
    if (data.length === 0) return '';
    
    const headers = Object.keys(data[0]);
    const csvContent = [
      headers.join(','),
      ...data.map(row => headers.map(header => row[header]).join(','))
    ].join('\n');
    
    return csvContent;
  }

  function downloadCSV(content: string, filename: string) {
    const blob = new Blob([content], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', filename);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  /**
   * 获取每个数据规模下指定性能指标的最佳结果
   * @param metric 性能指标字段名，如 'avgTime'
   * @returns 每个数据规模下的最佳结果对象
   */
  function getBestPerformance(metric: string) {
    if (results.length === 0) return null;
    
    const grouped = results.reduce((acc, result) => {
      if (!acc[result.dataSize]) acc[result.dataSize] = [];
      acc[result.dataSize].push(result);
      return acc;
    }, {});

    // 添加索引签名，允许用字符串索引
    const best: { [key: string]: any } = {};
    Object.keys(grouped).forEach(size => {
      const sizeResults = grouped[size];
      const bestResult = sizeResults.reduce(
        (best: any, current: any) => 
          parseFloat(current[metric]) < parseFloat(best[metric]) ? current : best
      );
      best[size] = bestResult;
    });

    return best;
  }
</script>

<div class="benchmark-panel">
  <div class="panel-header">
    <h3>算法性能评测</h3>
    <div class="header-actions">
      {#if results.length > 0}
        <button class="action-btn" on:click={exportResults}>
          导出结果
        </button>
        <button class="action-btn secondary" on:click={clearResults}>
          清除结果
        </button>
      {/if}
    </div>
  </div>

  <div class="benchmark-content">
    <!-- 配置面板 -->
    <div class="config-panel">
      <h4>测试配置</h4>
      
      <!-- 算法选择 -->
      <div class="config-section">
        <span class="section-label">选择算法:</span>
        <div class="algorithm-grid">
          {#each availableAlgorithms as algorithm}
            <label class="algorithm-checkbox">
              <input
                type="checkbox"
                value={algorithm.id}
                checked={selectedAlgorithms.includes(algorithm.id)}
                on:change={() => toggleAlgorithmSelection(algorithm.id)}
              />
              <span>{algorithm.name}</span>
            </label>
          {/each}
        </div>
      </div>

      <!-- 数据配置 -->
      <div class="config-section">
        <label for={inputIds.dataType}>数据类型:</label>
        <select id={inputIds.dataType} bind:value={dataType}>
          {#each Object.values(DATA_TYPES) as type}
            <option value={type}>{DATA_TYPE_NAMES[type]}</option>
          {/each}
        </select>
      </div>

      <div class="config-section">
        <label for={inputIds.dataPattern}>数据模式:</label>
        <select id={inputIds.dataPattern} bind:value={dataPattern}>
          {#each Object.values(DATA_PATTERNS) as pattern}
            <option value={pattern}>{PATTERN_NAMES[pattern]}</option>
          {/each}
        </select>
      </div>

      <!-- 数据规模 -->
      <div class="config-section">
        <span class="section-label">数据规模:</span>
        <div class="data-sizes">
          {#each dataSizes as size, index}
            <div class="size-item">
              <span>{size}</span>
              <button class="remove-btn" on:click={() => removeDataSize(index)}>×</button>
            </div>
          {/each}
          <button class="add-btn" on:click={addDataSize}>+ 添加</button>
        </div>
      </div>

      <!-- 测试次数 -->
      <div class="config-section">
        <label for={inputIds.testCount}>测试次数:</label>
        <input
          id={inputIds.testCount}
          type="number"
          bind:value={testCount}
          min="1"
          max="10"
        />
      </div>

      <!-- 运行按钮 -->
      <div class="config-section">
        <button 
          class="run-btn"
          disabled={isRunning || selectedAlgorithms.length === 0}
          on:click={runBenchmark}
        >
          {#if isRunning}
            运行中...
          {:else}
            开始测试
          {/if}
        </button>
      </div>
    </div>

    <!-- 结果面板 -->
    <div class="results-panel">
      <h4>测试结果</h4>
      
      {#if isRunning}
        <div class="loading">
          <div class="loading-spinner"></div>
          <p>正在运行性能测试...</p>
        </div>
      {:else if results.length === 0}
        <div class="empty-results">
          <p>暂无测试结果</p>
          <p>请配置测试参数并点击"开始测试"</p>
        </div>
      {:else}
        <div class="results-table">
          <table>
            <thead>
              <tr>
                <th>算法</th>
                <th>数据规模</th>
                <th>平均时间(ms)</th>
                <th>最短时间(ms)</th>
                <th>最长时间(ms)</th>
                <th>内存使用(KB)</th>
                <th>操作次数</th>
              </tr>
            </thead>
            <tbody>
              {#each results as result}
                <tr>
                  <td>{result.algorithmName}</td>
                  <td>{result.dataSize}</td>
                  <td>{result.avgTime}</td>
                  <td>{result.minTime}</td>
                  <td>{result.maxTime}</td>
                  <td>{result.memoryUsage}</td>
                  <td>{result.operations}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>

        <!-- 简单的性能对比 -->
        <div class="performance-summary">
          <h5>性能总结</h5>
          <p>基于当前测试结果的简单分析：</p>
          <ul>
            <li>共测试了 {selectedAlgorithms.length} 个算法</li>
            <li>测试了 {dataSizes.length} 种数据规模</li>
            <li>总共进行了 {results.length} 次测试</li>
          </ul>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .benchmark-panel {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border-light);
    border-radius: 8px;
    min-height: 600px;
    color: var(--color-text-primary);
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 1rem;
    border-bottom: 1px solid var(--color-border);
  }

  .panel-header h3 {
    margin: 0;
    color: var(--color-text-primary);
  }

  .header-actions {
    display: flex;
    gap: 0.5rem;
  }

  .action-btn {
    padding: 0.5rem 1rem;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    background: var(--color-surface-elevated);
    color: var(--color-text-primary);
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn:hover {
    background: var(--color-primary-50);
    border-color: var(--color-primary-300);
    color: var(--color-primary-700);
  }

  .action-btn.secondary {
    color: var(--color-error-500);
    border-color: var(--color-error-50);
  }

  .benchmark-content {
    display: grid;
    grid-template-columns: 350px 1fr;
    gap: 2rem;
  }

  .config-panel,
  .results-panel {
    background: var(--color-surface-elevated);
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: var(--shadow-sm);
    border: 1px solid var(--color-border-light);
  }

  .config-panel h4,
  .results-panel h4 {
    margin: 0 0 1rem 0;
    color: var(--color-text-primary);
  }

  .config-section {
    margin-bottom: 1rem;
  }

  .config-section label,
  .config-section .section-label {
    display: block;
    margin-bottom: 0.5rem;
    color: var(--color-text-primary);
    font-weight: 500;
  }

  .algorithm-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 0.5rem;
  }

  .algorithm-checkbox {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border-light);
    border-radius: 4px;
    cursor: pointer;
  }

  .algorithm-checkbox:hover {
    background: var(--color-primary-50);
    border-color: var(--color-primary-300);
  }

  .data-sizes {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .size-item {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.25rem 0.5rem;
    background: var(--color-border-light);
    border-radius: 4px;
    font-size: 0.875rem;
    color: var(--color-text-primary);
  }

  .remove-btn {
    background: none;
    border: none;
    color: var(--color-error-500);
    cursor: pointer;
    font-size: 1rem;
    line-height: 1;
  }

  .add-btn {
    padding: 0.25rem 0.5rem;
    background: var(--color-primary-500);
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 0.875rem;
    cursor: pointer;
  }

  .run-btn {
    width: 100%;
    padding: 0.75rem;
    background: var(--color-success-500);
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s;
  }

  .run-btn:hover:not(:disabled) {
    background: var(--color-success-600);
  }

  .run-btn:disabled {
    background: var(--color-border);
    color: var(--color-text-muted);
    cursor: not-allowed;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 3rem;
    color: var(--color-text-secondary);
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 4px solid var(--color-border-light);
    border-top: 4px solid var(--color-primary-500);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .empty-results {
    text-align: center;
    padding: 3rem;
    color: var(--color-text-muted);
  }

  .results-table {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.875rem;
  }

  th,
  td {
    padding: 0.75rem;
    text-align: left;
    border-bottom: 1px solid #e5e7eb;
  }

  th {
    background: #f9fafb;
    font-weight: 600;
    color: #374151;
  }

  .performance-summary {
    margin-top: 2rem;
    padding: 1rem;
    background: #f0f9ff;
    border-radius: 6px;
    border-left: 4px solid #3b82f6;
  }

  .performance-summary h5 {
    margin: 0 0 0.5rem 0;
    color: #1f2937;
  }

  .performance-summary p {
    margin: 0 0 0.5rem 0;
    color: #4b5563;
  }

  .performance-summary ul {
    margin: 0;
    padding-left: 1.5rem;
    color: #4b5563;
  }

  @media (max-width: 1024px) {
    .benchmark-content {
      grid-template-columns: 1fr;
    }
  }
</style>
