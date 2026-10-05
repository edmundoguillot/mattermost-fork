"""Prepare a clean runtime context; never copy development config or patient data."""
import json
import shutil
from pathlib import Path

root = Path(__file__).resolve().parents[2]
context = root / '.local-dev/docker-arm64'
package = context / 'mattermost'
# A development watcher can overwrite shared packages during a production build.
for artifact in (root / 'webapp/platform/shared/dist').rglob('*.js'):
    if 'jsxDEV' in artifact.read_text():
        raise RuntimeError(f'Development JSX in {artifact}; stop webapp watchers and rebuild production assets')
config_path = package / 'config/config.json'
config = json.loads(config_path.read_text())
config['ServiceSettings'].update(SiteURL='', ListenAddress=':8065', EnableDesktopLandingPage=False, EnableSecurityFixAlert=False, EnableEmailInvitations=True)
config['SqlSettings'].update(DriverName='postgres', DataSource='')
config['FileSettings'].update(Directory='/mattermost/data')
config['LogSettings'].update(FileLocation='/mattermost/logs/mattermost.log')
config['TeamSettings'].update(SiteName='Medikids', EnableCustomBrand=True, CustomBrandText='Comunicación interna de Medikids', CustomDescriptionText='', EnableUserCreation=True, EnableOpenServer=False)
config['LocalizationSettings'].update(DefaultServerLocale='es', DefaultClientLocale='es')
config['EmailSettings'].update(EnableSignUpWithEmail=True, SendPushNotifications=False, PushNotificationServer='')
config['PluginSettings']['Enable'] = False
config['SupportSettings'].update(AboutLink='', HelpLink='', PrivacyPolicyLink='', TermsOfServiceLink='', ReportAProblemLink='')
config_path.write_text(json.dumps(config, indent=2) + '\n')
for folder in ('fonts', 'templates', 'i18n'):
    shutil.copytree(root / 'server' / folder, package / folder, dirs_exist_ok=True)
shutil.rmtree(package / 'client', ignore_errors=True)
shutil.copytree(root / 'webapp/channels/dist', package / 'client')
(package / 'branding').mkdir(exist_ok=True)
shutil.copy2(root / 'logo.png', package / 'branding/logo.png')
for filename in ('LICENSE.txt', 'LICENSE.enterprise', 'NOTICE.txt'):
    shutil.copy2(root / filename, package / filename)
shutil.copy2(root / 'server/build/MIT-COMPILED-LICENSE.md', package)
shutil.copy2(root / 'deploy/docker/entrypoint.sh', context)
