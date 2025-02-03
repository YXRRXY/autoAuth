import os
import sys
import fileinput
import codecs
import yaml

def create_config():
    """创建默认配置文件"""
    config = {
        'server': {
            'port': 8080,
            'mode': 'debug'
        },
        'database': {
            'host': 'localhost',
            'port': 3306,
            'username': 'root',
            'password': '',  # 需要用户手动设置
            'dbname': 'autoauth'
        },
        'jwt': {
            'secret': ''  # 程序会自动生成
        }
    }

    # 确保配置目录存在
    os.makedirs('configs', exist_ok=True)
    
    # 写入配置文件
    config_path = 'configs/config.yaml'
    if not os.path.exists(config_path):
        with open(config_path, 'w', encoding='utf-8') as f:
            yaml.dump(config, f, default_flow_style=False, allow_unicode=True)
        print(f"已创建配置文件: {config_path}")
    else:
        print(f"配置文件已存在: {config_path}")

def replace_in_file(file_path, old_module, new_module):
    """使用 utf-8 编码读取和写入文件"""
    try:
        with codecs.open(file_path, 'r', encoding='utf-8') as file:
            content = file.read()
        
        content = content.replace(old_module, new_module)
        
        with codecs.open(file_path, 'w', encoding='utf-8') as file:
            file.write(content)
    except Exception as e:
        print(f"处理文件 {file_path} 时出错: {str(e)}")
        raise

def replace_module_name(old_module, new_module):
    """替换所有文件中的模块名"""
    # 替换go.mod文件中的模块名
    if os.path.exists('go.mod'):
        replace_in_file('go.mod', old_module, new_module)
        print(f"已更新 go.mod 文件")
    else:
        print("警告: 未找到 go.mod 文件")
    
    # 遍历所有.go文件并替换导入路径
    files_updated = 0
    for root, _, files in os.walk('.'):
        for file in files:
            if file.endswith('.go'):
                file_path = os.path.join(root, file)
                try:
                    replace_in_file(file_path, old_module, new_module)
                    files_updated += 1
                except Exception as e:
                    print(f"处理文件 {file_path} 时出错: {str(e)}")
    print(f"已更新 {files_updated} 个 .go 文件")

def check_go_installation():
    """检查 Go 环境"""
    try:
        import subprocess
        result = subprocess.run(['go', 'version'], capture_output=True, text=True)
        if result.returncode == 0:
            print(f"检测到 Go 环境: {result.stdout.strip()}")
            return True
        return False
    except FileNotFoundError:
        print("未检测到 Go 环境，请先安装 Go")
        return False

def init_project():
    """初始化项目"""
    # 创建项目目录结构
    directories = [
        'cmd/api',
        'configs',
        'internal/api/handlers',
        'internal/api/middleware',
        'internal/config',
        'internal/dal/model',
        'internal/dal/query',
        'internal/service',
        'pkg/utils',
    ]
    
    for directory in directories:
        os.makedirs(directory, exist_ok=True)
        print(f"创建目录: {directory}")

def main():
    if len(sys.argv) != 2:
        print("使用方法: python script.py <新模块名>")
        print("例如: python script.py github.com/YXRRXY/autoAuth/backend")
        sys.exit(1)

    # 检查 Go 环境
    if not check_go_installation():
        sys.exit(1)

    old_module = "github.com/YXRRXY/autoAuth"
    new_module = sys.argv[1]
    
    try:
        print("开始初始化项目...")
        init_project()
        
        print("创建配置文件...")
        create_config()
        
        print(f"开始替换模块名...")
        print(f"从: {old_module}")
        print(f"到: {new_module}")
        replace_module_name(old_module, new_module)
        
        print("\n初始化完成!")
        print("\n后续步骤:")
        print("1. 修改 configs/config.yaml 中的数据库配置")
        print("2. 运行 'go mod tidy' 更新依赖")
        print("3. 运行 'go run cmd/api/main.go' 启动服务")
        
    except Exception as e:
        print(f"发生错误: {str(e)}")
        sys.exit(1)

if __name__ == "__main__":
    main() 