<script lang="ts">
  import { 
    currentStep, 
    totalSteps, 
    isPlaying, 
    isPaused,
    playSpeed,
    progress,
    controlsState,
    visualizationActions,
    visualizationUtils
  } from '../../stores/visualization';

  // 组件属性
  export let showSpeedControl = true;
  export let showProgressBar = true;
  export let showStepInfo = true;
  export let compact = false;

  // 本地状态
  let speedValue = $playSpeed;

  // 响应式语句
  $: {
    if (speedValue !== $playSpeed) {
      visualizationActions.setSpeed(speedValue);
    }
  }

  // 事件处理
  function handlePlay() {
    visualizationActions.play();
  }

  function handlePause() {
    visualizationActions.pause();
  }

  function handleStop() {
    visualizationActions.stop();
  }

  function handleNext() {
    visualizationActions.nextStep();
  }

  function handlePrev() {
    visualizationActions.prevStep();
  }

  function handleReset() {
    visualizationActions.reset();
  }

  /**
   * 处理进度条点击，按照点击位置跳转到目标步骤。
   */
  function handleProgressClick(event: MouseEvent) {
    if (!$totalSteps) return;

    const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
    const clickX = event.clientX - rect.left;
    const percentage = clickX / rect.width;
    const targetStep = Math.floor(percentage * ($totalSteps - 1));
    
    visualizationActions.goToStep(targetStep);
  }

  function handleSpeedChange() {
    visualizationActions.setSpeed(speedValue);
  }

  /**
   * 处理进度条的键盘操作，支持常见的方向键与 Home/End 快捷键。
   */
  function handleProgressKey(event: KeyboardEvent) {
    if (!$totalSteps) {
      return;
    }

    let targetStep = $currentStep;

    switch (event.key) {
      case 'ArrowLeft':
      case 'ArrowDown':
        targetStep = Math.max(0, $currentStep - 1);
        break;
      case 'ArrowRight':
      case 'ArrowUp':
        targetStep = Math.min($totalSteps - 1, $currentStep + 1);
        break;
      case 'Home':
        targetStep = 0;
        break;
      case 'End':
        targetStep = $totalSteps - 1;
        break;
      default:
        return;
    }

    event.preventDefault();
    visualizationActions.goToStep(targetStep);
  }

  // 格式化速度显示
  function formatSpeed(speed: number): string {
    if (speed >= 1000) {
      return `${(speed / 1000).toFixed(1)}s`;
    }
    return `${speed}ms`;
  }

  // 获取播放按钮图标
  function getPlayIcon(): string {
    if ($isPlaying) return '⏸️';
    if ($isPaused) return '▶️';
    return '▶️';
  }

  // 获取播放按钮文本
  function getPlayText(): string {
    if ($isPlaying) return '暂停';
    if ($isPaused) return '继续';
    return '播放';
  }
</script>

<div class="visualization-controls" class:compact>
  <!-- 主要控制按钮 -->
  <div class="main-controls">
    <button
      class="control-btn play-btn"
      class:playing={$isPlaying}
      disabled={$isPlaying ? !$controlsState.canPause : !$controlsState.canPlay}
      on:click={$isPlaying ? handlePause : handlePlay}
      title={getPlayText()}
    >
      <span class="icon">{getPlayIcon()}</span>
      {#if !compact}
        <span class="text">{getPlayText()}</span>
      {/if}
    </button>

    <button
      class="control-btn stop-btn"
      disabled={!$isPlaying && !$isPaused}
      on:click={handleStop}
      title="停止"
    >
      <span class="icon">⏹️</span>
      {#if !compact}
        <span class="text">停止</span>
      {/if}
    </button>

    <div class="step-controls">
      <button
        class="control-btn step-btn"
        disabled={!$controlsState.canPrev}
        on:click={handlePrev}
        title="上一步"
      >
        <span class="icon">⏮️</span>
      </button>

      <button
        class="control-btn step-btn"
        disabled={!$controlsState.canNext}
        on:click={handleNext}
        title="下一步"
      >
        <span class="icon">⏭️</span>
      </button>
    </div>

    <button
      class="control-btn reset-btn"
      disabled={!$controlsState.canReset}
      on:click={handleReset}
      title="重置"
    >
      <span class="icon">🔄</span>
      {#if !compact}
        <span class="text">重置</span>
      {/if}
    </button>
  </div>

  <!-- 进度条 -->
  {#if showProgressBar && $totalSteps > 0}
    <div class="progress-section">
      <div 
        class="progress-bar"
        on:click={handleProgressClick}
        on:keydown={handleProgressKey}
        role="slider"
        tabindex="0"
        aria-valuemin="0"
        aria-valuemax={$totalSteps - 1}
        aria-valuenow={$currentStep}
      >
        <div 
          class="progress-fill"
          style="width: {$progress}%"
        ></div>
        <div 
          class="progress-thumb"
          style="left: {$progress}%"
        ></div>
      </div>
    </div>
  {/if}

  <!-- 步骤信息 -->
  {#if showStepInfo && $totalSteps > 0}
    <div class="step-info">
      <span class="step-text">
        步骤 {$currentStep + 1} / {$totalSteps}
      </span>
      <span class="progress-text">
        ({$progress.toFixed(1)}%)
      </span>
    </div>
  {/if}

  <!-- 速度控制 -->
  {#if showSpeedControl}
    <div class="speed-control">
      <label for="speed-slider">播放速度:</label>
      <input
        id="speed-slider"
        type="range"
        min="100"
        max="3000"
        step="100"
        bind:value={speedValue}
        on:change={handleSpeedChange}
      />
      <span class="speed-value">{formatSpeed(speedValue)}</span>
    </div>
  {/if}
</div>

<style>
  .visualization-controls {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border-light);
    border-radius: 8px;
    box-shadow: var(--shadow-sm);
    color: var(--color-text-primary);
  }

  .visualization-controls.compact {
    padding: 0.5rem;
    gap: 0.5rem;
  }

  .main-controls {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    justify-content: center;
  }

  .control-btn {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.5rem 1rem;
    background: var(--color-surface-elevated);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    font-size: 0.875rem;
    color: var(--color-text-primary);
    cursor: pointer;
    transition: all 0.2s;
  }

  .control-btn:hover:not(:disabled) {
    background: var(--color-primary-50);
    border-color: var(--color-primary-300);
    color: var(--color-primary-700);
  }

  .control-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    color: var(--color-text-muted);
  }

  .play-btn.playing {
    background: var(--color-warning-50);
    border-color: var(--color-warning-500);
    color: var(--color-warning-500);
    font-weight: 500;
  }

  .step-controls {
    display: flex;
    gap: 0.25rem;
  }

  .step-btn {
    padding: 0.5rem;
    min-width: auto;
  }

  .icon {
    font-size: 1rem;
  }

  .text {
    font-size: 0.875rem;
    font-weight: 500;
    color: inherit;
  }

  .progress-section {
    width: 100%;
  }

  .progress-bar {
    position: relative;
    width: 100%;
    height: 8px;
    background: var(--color-border-light);
    border-radius: 4px;
    cursor: pointer;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: var(--color-primary-500);
    border-radius: 4px;
    transition: width 0.2s ease;
  }

  .progress-thumb {
    position: absolute;
    top: -4px;
    width: 16px;
    height: 16px;
    background: var(--color-primary-500);
    border: 2px solid var(--color-surface);
    border-radius: 50%;
    transform: translateX(-50%);
    box-shadow: var(--shadow-md);
    transition: left 0.2s ease;
  }

  .step-info {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    color: var(--color-text-secondary);
  }

  .step-text {
    font-weight: 500;
    color: var(--color-text-primary);
  }

  .speed-control {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
  }

  .speed-control label {
    color: var(--color-text-primary);
    white-space: nowrap;
    font-weight: 500;
  }

  .speed-control input[type="range"] {
    flex: 1;
    min-width: 100px;
  }

  .speed-value {
    color: var(--color-text-secondary);
    min-width: 50px;
    text-align: right;
    font-weight: 500;
  }

  @media (max-width: 640px) {
    .main-controls {
      flex-wrap: wrap;
    }

    .control-btn .text {
      display: none;
    }

    .speed-control {
      flex-direction: column;
      align-items: stretch;
      gap: 0.25rem;
    }
  }

  /* 深色模式特殊优化 */
  @media (prefers-color-scheme: dark) {
    .control-btn:hover:not(:disabled) {
      background: var(--color-primary-800);
      border-color: var(--color-primary-600);
      color: var(--color-primary-100);
    }

    .play-btn.playing {
      background: var(--color-warning-800);
      border-color: var(--color-warning-600);
      color: var(--color-warning-200);
    }

    .icon {
      filter: brightness(1.2);
    }
  }
</style>
