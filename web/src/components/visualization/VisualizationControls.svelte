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

  function handleProgressClick(event: MouseEvent) {
    if (!$totalSteps) return;

    const rect = (event.target as HTMLElement).getBoundingClientRect();
    const clickX = event.clientX - rect.left;
    const percentage = clickX / rect.width;
    const targetStep = Math.floor(percentage * ($totalSteps - 1));
    
    visualizationActions.goToStep(targetStep);
  }

  function handleSpeedChange() {
    visualizationActions.setSpeed(speedValue);
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
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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
    background: #f3f4f6;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .control-btn:hover:not(:disabled) {
    background: #e5e7eb;
    border-color: #9ca3af;
  }

  .control-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .play-btn.playing {
    background: #fef3c7;
    border-color: #f59e0b;
    color: #92400e;
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
  }

  .progress-section {
    width: 100%;
  }

  .progress-bar {
    position: relative;
    width: 100%;
    height: 8px;
    background: #e5e7eb;
    border-radius: 4px;
    cursor: pointer;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: #3b82f6;
    border-radius: 4px;
    transition: width 0.2s ease;
  }

  .progress-thumb {
    position: absolute;
    top: -4px;
    width: 16px;
    height: 16px;
    background: #3b82f6;
    border: 2px solid white;
    border-radius: 50%;
    transform: translateX(-50%);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    transition: left 0.2s ease;
  }

  .step-info {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    color: #6b7280;
  }

  .step-text {
    font-weight: 500;
    color: #374151;
  }

  .speed-control {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
  }

  .speed-control label {
    color: #374151;
    white-space: nowrap;
  }

  .speed-control input[type="range"] {
    flex: 1;
    min-width: 100px;
  }

  .speed-value {
    color: #6b7280;
    min-width: 50px;
    text-align: right;
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
</style>
