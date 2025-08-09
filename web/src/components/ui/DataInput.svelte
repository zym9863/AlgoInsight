<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { DataType, DataPattern } from '../../types/data';
  import { DATA_TYPES, DATA_PATTERNS, DATA_TYPE_NAMES, PATTERN_NAMES } from '../../types/data';

  // 事件分发器
  const dispatch = createEventDispatcher<{
    dataChange: any;
    generate: { dataType: DataType; size: number; pattern: DataPattern };
  }>();

  // 组件属性
  export let value: any = null;
  export let dataType: DataType = DATA_TYPES.ARRAY;
  export let allowGeneration = true;
  export let allowCustomInput = true;
  export let maxSize = 100;

  // 本地状态
  let inputMode: 'custom' | 'generate' | 'preset' = 'custom';
  let customInput = '';
  let generateSize = 10;
  let generatePattern: DataPattern = DATA_PATTERNS.RANDOM;
  let parseError = '';

  // 响应式语句
  $: {
    if (inputMode === 'custom' && customInput) {
      parseCustomInput();
    }
  }

  // 解析自定义输入
  function parseCustomInput() {
    parseError = '';
    
    try {
      if (!customInput.trim()) {
        value = null;
        dispatch('dataChange', null);
        return;
      }

      let parsed;
      
      if (dataType === DATA_TYPES.ARRAY) {
        // 尝试解析为数组
        if (customInput.startsWith('[') && customInput.endsWith(']')) {
          parsed = JSON.parse(customInput);
        } else {
          // 逗号分隔的值
          parsed = customInput.split(',').map(item => {
            const trimmed = item.trim();
            // 尝试解析为数字
            const num = Number(trimmed);
            return isNaN(num) ? trimmed : num;
          });
        }
      } else {
        // 其他数据类型尝试JSON解析
        parsed = JSON.parse(customInput);
      }

      value = parsed;
      dispatch('dataChange', parsed);
    } catch (error) {
      parseError = '输入格式错误，请检查数据格式';
      console.error('Parse error:', error);
    }
  }

  // 生成数据
  function generateData() {
    if (generateSize <= 0 || generateSize > maxSize) {
      parseError = `数据大小必须在 1 到 ${maxSize} 之间`;
      return;
    }

    parseError = '';
    dispatch('generate', {
      dataType,
      size: generateSize,
      pattern: generatePattern
    });
  }

  // 设置预设数据
  function setPresetData(data: any) {
    value = data;
    customInput = JSON.stringify(data);
    inputMode = 'custom';
    dispatch('dataChange', data);
  }

  // 获取数据预览
  function getDataPreview(data: any): string {
    if (!data) return '无数据';
    
    if (Array.isArray(data)) {
      if (data.length <= 10) {
        return `[${data.join(', ')}]`;
      } else {
        return `[${data.slice(0, 5).join(', ')}, ..., ${data.slice(-2).join(', ')}] (${data.length}个元素)`;
      }
    }
    
    return JSON.stringify(data).substring(0, 100) + (JSON.stringify(data).length > 100 ? '...' : '');
  }

  // 预设数据示例
  // 增加索引签名，允许用 string 作为 key，解决 TS 报错
  const presetData: { [key: string]: { name: string; data: number[] }[] } = {
    [DATA_TYPES.ARRAY]: [
      { name: '小型随机数组', data: [64, 34, 25, 12, 22, 11, 90, 5, 77, 30] },
      { name: '已排序数组', data: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10] },
      { name: '逆序数组', data: [10, 9, 8, 7, 6, 5, 4, 3, 2, 1] },
      { name: '重复元素', data: [5, 2, 8, 2, 9, 1, 5, 5, 2, 8] }
    ]
  };
</script>

<div class="data-input">
  <div class="header">
    <h4>数据输入</h4>
    <div class="data-type-selector">
      <label>数据类型:</label>
      <select bind:value={dataType}>
        {#each Object.values(DATA_TYPES) as type}
          <option value={type}>{DATA_TYPE_NAMES[type]}</option>
        {/each}
      </select>
    </div>
  </div>

  <!-- 输入模式选择 -->
  <div class="mode-tabs">
    {#if allowCustomInput}
      <button
        class="mode-tab"
        class:active={inputMode === 'custom'}
        on:click={() => inputMode = 'custom'}
      >
        自定义输入
      </button>
    {/if}
    
    {#if allowGeneration}
      <button
        class="mode-tab"
        class:active={inputMode === 'generate'}
        on:click={() => inputMode = 'generate'}
      >
        生成数据
      </button>
    {/if}
    
    <button
      class="mode-tab"
      class:active={inputMode === 'preset'}
      on:click={() => inputMode = 'preset'}
    >
      预设数据
    </button>
  </div>

  <!-- 自定义输入 -->
  {#if inputMode === 'custom' && allowCustomInput}
    <div class="custom-input">
      <label>
        {#if dataType === DATA_TYPES.ARRAY}
          输入数组 (支持JSON格式或逗号分隔):
        {:else}
          输入JSON格式数据:
        {/if}
      </label>
      <textarea
        bind:value={customInput}
        placeholder={dataType === DATA_TYPES.ARRAY 
          ? "例如: [1, 2, 3, 4, 5] 或 1, 2, 3, 4, 5" 
          : "输入JSON格式数据"}
        rows="4"
      ></textarea>
      
      {#if parseError}
        <div class="error">{parseError}</div>
      {/if}
    </div>
  {/if}

  <!-- 生成数据 -->
  {#if inputMode === 'generate' && allowGeneration}
    <div class="generate-data">
      <div class="generate-controls">
        <div class="control-group">
          <label>数据大小:</label>
          <input
            type="number"
            bind:value={generateSize}
            min="1"
            max={maxSize}
          />
        </div>
        
        <div class="control-group">
          <label>数据模式:</label>
          <select bind:value={generatePattern}>
            {#each Object.values(DATA_PATTERNS) as pattern}
              <option value={pattern}>{PATTERN_NAMES[pattern]}</option>
            {/each}
          </select>
        </div>
      </div>
      
      <button class="generate-btn" on:click={generateData}>
        生成数据
      </button>
      
      {#if parseError}
        <div class="error">{parseError}</div>
      {/if}
    </div>
  {/if}

  <!-- 预设数据 -->
  {#if inputMode === 'preset'}
    <div class="preset-data">
      {#if dataType in presetData}
        <div class="preset-list">
          {#each presetData[dataType] as preset}
            <button
              class="preset-item"
              on:click={() => setPresetData(preset.data)}
            >
              <div class="preset-name">{preset.name}</div>
              <div class="preset-preview">{getDataPreview(preset.data)}</div>
            </button>
          {/each}
        </div>
      {:else}
        <div class="no-presets">
          暂无 {DATA_TYPE_NAMES[dataType]} 类型的预设数据
        </div>
      {/if}
    </div>
  {/if}

  <!-- 数据预览 -->
  {#if value}
    <div class="data-preview">
      <h5>当前数据预览:</h5>
      <div class="preview-content">
        {getDataPreview(value)}
      </div>
    </div>
  {/if}
</div>

<style>
  .data-input {
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

  .data-input::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, var(--color-primary-500), var(--color-primary-600));
  }

  .data-input:hover {
    box-shadow: var(--shadow-md);
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--spacing-sm);
  }

  .header h4 {
    margin: 0;
    color: var(--color-text-primary);
    font-weight: 700;
    font-size: 1.125rem;
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
  }

  .header h4::before {
    content: '📊';
    font-size: 1.25rem;
  }

  .data-type-selector {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    font-size: 0.8125rem;
  }

  .data-type-selector label {
    color: var(--color-text-secondary);
    font-weight: 500;
  }

  .data-type-selector select {
    padding: var(--spacing-xs) var(--spacing-sm);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    font-size: 0.8125rem;
    background: var(--color-surface);
    color: var(--color-text-primary);
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .data-type-selector select:focus {
    outline: none;
    border-color: var(--color-primary-500);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .mode-tabs {
    display: flex;
    gap: var(--spacing-xs);
    border-bottom: 1px solid var(--color-border-light);
    padding-bottom: var(--spacing-md);
    margin-bottom: var(--spacing-sm);
  }

  .mode-tab {
    padding: var(--spacing-sm) var(--spacing-md);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    font-size: 0.8125rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-fast);
    color: var(--color-text-secondary);
    position: relative;
    flex: 1;
    text-align: center;
  }

  .mode-tab::before {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
    border-radius: inherit;
    opacity: 0;
    transition: opacity var(--transition-fast);
  }

  .mode-tab:hover {
    background: var(--color-gray-50);
    border-color: var(--color-primary-300);
    color: var(--color-primary-700);
    transform: translateY(-1px);
  }

  .mode-tab.active {
    background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
    color: white;
    border-color: var(--color-primary-500);
    box-shadow: var(--shadow-sm);
  }

  .custom-input {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .custom-input label {
    display: block;
    color: var(--color-text-primary);
    font-size: 0.8125rem;
    font-weight: 600;
    margin-bottom: var(--spacing-xs);
  }

  .custom-input textarea {
    width: 100%;
    padding: var(--spacing-md);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    font-family: var(--font-family-mono);
    font-size: 0.8125rem;
    resize: vertical;
    background: var(--color-surface);
    color: var(--color-text-primary);
    transition: all var(--transition-fast);
    min-height: 100px;
  }

  .custom-input textarea::placeholder {
    color: var(--color-text-muted);
    font-style: italic;
  }

  .custom-input textarea:focus {
    outline: none;
    border-color: var(--color-primary-500);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    background: var(--color-surface-elevated);
  }

  .generate-data {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .generate-controls {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-md);
  }

  .control-group {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-xs);
  }

  .control-group label {
    color: var(--color-text-primary);
    font-size: 0.8125rem;
    font-weight: 600;
  }

  .control-group input,
  .control-group select {
    padding: var(--spacing-sm);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    font-size: 0.8125rem;
    background: var(--color-surface);
    color: var(--color-text-primary);
    transition: all var(--transition-fast);
  }

  .control-group input:focus,
  .control-group select:focus {
    outline: none;
    border-color: var(--color-primary-500);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .generate-btn {
    padding: var(--spacing-md) var(--spacing-lg);
    background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
    color: white;
    border: none;
    border-radius: var(--radius-md);
    font-size: 0.8125rem;
    font-weight: 600;
    cursor: pointer;
    transition: all var(--transition-fast);
    align-self: flex-start;
    box-shadow: var(--shadow-sm);
    position: relative;
    overflow: hidden;
  }

  .generate-btn::before {
    content: '⚡';
    margin-right: var(--spacing-xs);
  }

  .generate-btn:hover {
    background: linear-gradient(135deg, var(--color-primary-600), var(--color-primary-700));
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }

  .generate-btn:active {
    transform: translateY(0);
    box-shadow: var(--shadow-sm);
  }

  .preset-list {
    display: grid;
    gap: var(--spacing-sm);
  }

  .preset-item {
    padding: var(--spacing-md);
    background: var(--color-surface);
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-md);
    text-align: left;
    cursor: pointer;
    transition: all var(--transition-normal);
    position: relative;
    overflow: hidden;
  }

  .preset-item::before {
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

  .preset-item:hover {
    background: var(--color-gray-50);
    border-color: var(--color-primary-300);
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
  }

  .preset-item:hover::before {
    transform: scaleY(1);
  }

  .preset-name {
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: var(--spacing-xs);
    font-size: 0.875rem;
  }

  .preset-preview {
    font-size: 0.75rem;
    color: var(--color-text-muted);
    font-family: var(--font-family-mono);
    line-height: 1.4;
    background: var(--color-gray-50);
    padding: var(--spacing-xs);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-light);
  }

  .no-presets {
    text-align: center;
    color: var(--color-text-muted);
    padding: var(--spacing-2xl);
    background: var(--color-gray-50);
    border-radius: var(--radius-md);
    border: 2px dashed var(--color-border);
  }

  .no-presets::before {
    content: '📁';
    display: block;
    font-size: 2rem;
    margin-bottom: var(--spacing-sm);
    opacity: 0.5;
  }

  .data-preview {
    padding: var(--spacing-md);
    background: linear-gradient(135deg, var(--color-gray-50), rgba(59, 130, 246, 0.02));
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border-light);
    position: relative;
    overflow: hidden;
  }

  .data-preview::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, var(--color-success-500), var(--color-primary-500));
  }

  .data-preview h5 {
    margin: 0 0 var(--spacing-sm) 0;
    color: var(--color-text-primary);
    font-size: 0.8125rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
  }

  .data-preview h5::before {
    content: '👁️';
    font-size: 0.875rem;
  }

  .preview-content {
    font-family: var(--font-family-mono);
    font-size: 0.75rem;
    color: var(--color-text-secondary);
    word-break: break-all;
    background: var(--color-surface);
    padding: var(--spacing-sm);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-light);
    line-height: 1.5;
  }

  .error {
    color: var(--color-error-500);
    font-size: 0.8125rem;
    font-weight: 500;
    padding: var(--spacing-sm);
    background: var(--color-error-50);
    border: 1px solid rgba(239, 68, 68, 0.2);
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
  }

  .error::before {
    content: '⚠️';
    font-size: 0.875rem;
  }

  /* 响应式设计 */
  @media (max-width: 768px) {
    .data-input {
      padding: var(--spacing-md);
    }

    .header {
      flex-direction: column;
      align-items: flex-start;
      gap: var(--spacing-sm);
    }

    .mode-tabs {
      flex-direction: column;
    }

    .generate-controls {
      grid-template-columns: 1fr;
    }

    .mode-tab {
      text-align: center;
    }
  }

  /* 加载动画 */
  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .custom-input,
  .generate-data,
  .preset-data,
  .data-preview {
    animation: slideIn 0.3s ease-out;
  }

  /* 脉冲动画用于高亮重要元素 */
  @keyframes pulse {
    0%, 100% {
      opacity: 1;
    }
    50% {
      opacity: 0.8;
    }
  }

  .generate-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }
</style>
