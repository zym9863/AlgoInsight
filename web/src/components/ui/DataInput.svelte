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
    gap: 1rem;
    padding: 1rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .header h4 {
    margin: 0;
    color: #1f2937;
  }

  .data-type-selector {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
  }

  .data-type-selector label {
    color: #6b7280;
  }

  .data-type-selector select {
    padding: 0.25rem 0.5rem;
    border: 1px solid #d1d5db;
    border-radius: 4px;
    font-size: 0.875rem;
  }

  .mode-tabs {
    display: flex;
    gap: 0.5rem;
    border-bottom: 1px solid #e5e7eb;
    padding-bottom: 0.5rem;
  }

  .mode-tab {
    padding: 0.5rem 1rem;
    background: none;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .mode-tab:hover {
    background: #f3f4f6;
  }

  .mode-tab.active {
    background: #3b82f6;
    color: white;
    border-color: #3b82f6;
  }

  .custom-input label {
    display: block;
    margin-bottom: 0.5rem;
    color: #374151;
    font-size: 0.875rem;
  }

  .custom-input textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-family: 'Courier New', monospace;
    font-size: 0.875rem;
    resize: vertical;
  }

  .custom-input textarea:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .generate-data {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .generate-controls {
    display: flex;
    gap: 1rem;
  }

  .control-group {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .control-group label {
    color: #374151;
    font-size: 0.875rem;
  }

  .control-group input,
  .control-group select {
    padding: 0.5rem;
    border: 1px solid #d1d5db;
    border-radius: 4px;
    font-size: 0.875rem;
  }

  .generate-btn {
    padding: 0.75rem 1.5rem;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.875rem;
    cursor: pointer;
    transition: background 0.2s;
    align-self: flex-start;
  }

  .generate-btn:hover {
    background: #2563eb;
  }

  .preset-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .preset-item {
    padding: 0.75rem;
    background: #f9fafb;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    text-align: left;
    cursor: pointer;
    transition: all 0.2s;
  }

  .preset-item:hover {
    background: #f3f4f6;
    border-color: #3b82f6;
  }

  .preset-name {
    font-weight: 500;
    color: #1f2937;
    margin-bottom: 0.25rem;
  }

  .preset-preview {
    font-size: 0.75rem;
    color: #6b7280;
    font-family: 'Courier New', monospace;
  }

  .no-presets {
    text-align: center;
    color: #6b7280;
    padding: 2rem;
  }

  .data-preview {
    padding: 1rem;
    background: #f9fafb;
    border-radius: 6px;
    border: 1px solid #e5e7eb;
  }

  .data-preview h5 {
    margin: 0 0 0.5rem 0;
    color: #374151;
    font-size: 0.875rem;
  }

  .preview-content {
    font-family: 'Courier New', monospace;
    font-size: 0.75rem;
    color: #6b7280;
    word-break: break-all;
  }

  .error {
    color: #dc2626;
    font-size: 0.875rem;
    padding: 0.5rem;
    background: #fef2f2;
    border: 1px solid #fecaca;
    border-radius: 4px;
  }
</style>
