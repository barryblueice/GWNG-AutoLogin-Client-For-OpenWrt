import subprocess  
  
for arch in ["linux/386", "linux/amd64", "linux/arm", "linux/arm64", "linux/mips", "linux/mips64", "linux/mips64le", "linux/mipsle", "linux/ppc64", "linux/ppc64le", "linux/riscv64", "linux/s390x"]:  
    os_name, arch_name = arch.split("/")  
      
    ps_command = f'$env:CGO_ENABLED="0"; $env:GOOS="{os_name}"; $env:GOARCH="{arch_name}"; go build -o .\\build\gwng-{os_name}-{arch_name} gwng.go'  
      
    subprocess.run(["powershell", "-Command", ps_command], check=True, shell=True)
    
ps_command = f'$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="mipsle"; $env:GOMIPS="softfloat"; go build -o .\\build\gwng-linux-mipsle-mt7621 gwng.go'

subprocess.run(["powershell", "-Command", ps_command], check=True, shell=True)

# import os

# file_list = os.listdir(os.path.join(os.getcwd(),'release'))

# print (file_list)