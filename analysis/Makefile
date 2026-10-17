
## ai/init: 安装 AI 分析模块依赖（本地开发用）
ai/init:
	@echo "安装基础依赖..."
	@pip install -r analysis/requirements.txt
	@echo ""
	@echo "TFLite 支持需要额外安装（二选一）:"
	@echo "  Linux:        make ai/init/tflite"
	@echo "  macOS/Windows: make ai/init/tensorflow"
	@echo ""
	@echo "如果使用 .tflite 模型，请根据系统选择安装"

## ai/init/tflite: 安装 TFLite 支持（Linux 轻量版）
ai/init/tflite:
	@conda install tflite-runtime

## ai/init/tensorflow: 安装 TensorFlow（macOS/Windows 完整版）
ai/init/tensorflow:
	@ pip install tensorflow -i https://pypi.org/simple